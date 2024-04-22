package commands

import (
	"context"
	"crypto/x509"
	"encoding/pem"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"panel/app/models"
	"panel/internal"
	"panel/internal/services"
)

// CertRenew 证书续签
type CertRenew struct {
}

// Signature The name and signature of the console command.
func (receiver *CertRenew) Signature() string {
	return "panel:cert-renew"
}

// Description The console command description.
func (receiver *CertRenew) Description() string {
	ctx := context.Background()
	return facades.Lang(ctx).Get("commands.panel:cert-renew.description")
}

// Extend The console command extend.
func (receiver *CertRenew) Extend() command.Extend {
	return command.Extend{
		Category: "panel",
	}
}

// Handle Execute the console command.
func (receiver *CertRenew) Handle(console.Context) error {
	if internal.Status != internal.StatusNormal {
		return nil
	}

	var certs []models.Cert
	err := facades.Orm().Query().With("Website").With("User").With("DNS").Find(&certs)
	if err != nil {
		return err
	}

	for _, cert := range certs {
		if !cert.AutoRenew {
			continue
		}

		block, _ := pem.Decode([]byte(cert.Cert))
		if block != nil {
			data, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				continue
			}

			// 结束时间大于 7 天的证书不续签
			endTime := carbon.FromStdTime(data.NotAfter)
			if endTime.Gt(carbon.Now().AddDays(7)) {
				continue
			}
		}

		certService := services.NewCertImpl()
		_, err = certService.Renew(cert.ID)
		if err != nil {
			facades.Log().Tags("面板", "证书管理").With(map[string]any{
				"cert_id": cert.ID,
				"error":   err.Error(),
			}).Infof("证书续签失败")
		}
	}

	return nil
}
