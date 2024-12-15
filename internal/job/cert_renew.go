package job

import (
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	pkgcert "github.com/TheTNB/panel/pkg/cert"
)

// CertRenew 证书续签
type CertRenew struct {
	db       *gorm.DB
	log      *slog.Logger
	certRepo biz.CertRepo
}

func NewCertRenew(db *gorm.DB, log *slog.Logger, cert biz.CertRepo) *CertRenew {
	return &CertRenew{
		db:       db,
		log:      log,
		certRepo: cert,
	}
}

func (r *CertRenew) Run() {
	if app.Status != app.StatusNormal {
		return
	}

	var certs []biz.Cert
	if err := r.db.Preload("Website").Preload("Account").Preload("DNS").Find(&certs).Error; err != nil {
		r.log.Warn("获取证书失败", slog.Any("err", err))
		return
	}

	for _, cert := range certs {
		if cert.Type == "upload" || !cert.AutoRenew {
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

		_, err = r.certRepo.Renew(cert.ID)
		if err != nil {
			r.log.Warn("续签证书失败", slog.Any("err", err))
		}
	}
}
