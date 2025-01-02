package job

import (
	"log/slog"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"github.com/tnb-labs/panel/internal/app"
	"github.com/tnb-labs/panel/internal/biz"
	pkgcert "github.com/tnb-labs/panel/pkg/cert"
	"github.com/tnb-labs/panel/pkg/io"
	"github.com/tnb-labs/panel/pkg/shell"
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

	// 续签面板证书
	panelCert, err := io.Read(filepath.Join(app.Root, "panel/storage/cert.pem"))
	if err != nil {
		r.log.Warn("读取面板证书失败", slog.Any("err", err))
		return
	}
	decode, err := pkgcert.ParseCert(panelCert)
	if err != nil {
		r.log.Warn("解析面板证书失败", slog.Any("err", err))
		return
	}
	if decode.NotAfter.Sub(time.Now()).Hours() < 24*7 {
		_, err = shell.Exec("panel-cli https generate")
	}
}
