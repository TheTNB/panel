package job

import (
	"time"

	"go.uber.org/zap"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	pkgcert "github.com/TheTNB/panel/pkg/cert"
	"github.com/TheTNB/panel/pkg/types"
)

// CertRenew 证书续签
type CertRenew struct {
	cert biz.CertRepo
}

func NewCertRenew() *CertRenew {
	return &CertRenew{
		cert: data.NewCertRepo(),
	}
}

func (receiver *CertRenew) Run() {
	if types.Status != types.StatusNormal {
		return
	}

	var certs []biz.Cert
	if err := app.Orm.Preload("Website").Preload("Account").Preload("DNS").Find(&certs).Error; err != nil {
		app.Logger.Error("获取证书失败", zap.Error(err))
		return
	}

	for _, cert := range certs {
		if !cert.AutoRenew {
			continue
		}

		decode, err := pkgcert.ParseCert(cert.Cert)
		if err != nil {
			continue
		}

		// 结束时间大于 7 天的证书不续签
		now := time.Now()
		if decode.NotAfter.Sub(now).Hours() > 24*7 {
			continue
		}

		_, err = receiver.cert.Renew(cert.ID)
		if err != nil {
			app.Logger.Error("续签证书失败", zap.Error(err))
		}
	}
}
