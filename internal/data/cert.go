package data

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/acme"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type certRepo struct {
	client      *acme.Client
	websiteRepo biz.WebsiteRepo
}

func NewCertRepo() biz.CertRepo {
	return &certRepo{
		websiteRepo: NewWebsiteRepo(),
	}
}

func (r *certRepo) List(page, limit uint) ([]*biz.Cert, int64, error) {
	var certs []*biz.Cert
	var total int64
	err := app.Orm.Model(&biz.Cert{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&certs).Error
	return certs, total, err
}

func (r *certRepo) Get(id uint) (*biz.Cert, error) {
	cert := new(biz.Cert)
	err := app.Orm.Model(&biz.Cert{}).Where("id = ?", id).First(cert).Error
	return cert, err
}

func (r *certRepo) Create(req *request.CertCreate) (*biz.Cert, error) {
	cert := &biz.Cert{
		UserID:    req.UserID,
		WebsiteID: req.WebsiteID,
		DNSID:     req.DNSID,
		Type:      req.Type,
		Domains:   req.Domains,
		AutoRenew: req.AutoRenew,
	}
	if err := app.Orm.Create(cert).Error; err != nil {
		return nil, err
	}
	return cert, nil
}

func (r *certRepo) Update(req *request.CertUpdate) error {
	return app.Orm.Model(&biz.Cert{}).Where("id = ?", req.ID).Updates(&biz.Cert{
		UserID:    req.UserID,
		WebsiteID: req.WebsiteID,
		DNSID:     req.DNSID,
		Type:      req.Type,
		Domains:   req.Domains,
		AutoRenew: req.AutoRenew,
	}).Error
}

func (r *certRepo) Delete(id uint) error {
	return app.Orm.Model(&biz.Cert{}).Where("id = ?", id).Delete(&biz.Cert{}).Error
}

func (r *certRepo) ObtainAuto(id uint) (*acme.Certificate, error) {
	cert, err := r.Get(id)
	if err != nil {
		return nil, err
	}

	client, err := r.getClient(cert)
	if err != nil {
		return nil, err
	}

	if cert.DNS != nil {
		client.UseDns(acme.DnsType(cert.DNS.Type), cert.DNS.Data)
	} else {
		if cert.Website == nil {
			return nil, errors.New("该证书没有关联网站，无法自动签发")
		} else {
			for _, domain := range cert.Domains {
				if strings.Contains(domain, "*") {
					return nil, errors.New("通配符域名无法使用 HTTP 验证")
				}
			}
			conf := fmt.Sprintf("%s/server/vhost/acme/%s.conf", app.Root, cert.Website.Name)
			client.UseHTTP(conf, cert.Website.Path)
		}
	}

	ssl, err := client.ObtainSSL(context.Background(), cert.Domains, acme.KeyType(cert.Type))
	if err != nil {
		return nil, err
	}

	cert.CertURL = ssl.URL
	cert.Cert = string(ssl.ChainPEM)
	cert.Key = string(ssl.PrivateKey)
	if err = app.Orm.Save(cert).Error; err != nil {
		return nil, err
	}

	if cert.Website != nil {
		return &ssl, r.Deploy(cert.ID, cert.WebsiteID)
	}

	return &ssl, nil
}

func (r *certRepo) ObtainManual(id uint) (*acme.Certificate, error) {
	cert, err := r.Get(id)
	if err != nil {
		return nil, err
	}

	if r.client == nil {
		return nil, errors.New("请重新获取 DNS 解析记录")
	}

	ssl, err := r.client.ObtainSSLManual()
	if err != nil {
		return nil, err
	}

	cert.CertURL = ssl.URL
	cert.Cert = string(ssl.ChainPEM)
	cert.Key = string(ssl.PrivateKey)
	if err = app.Orm.Save(cert).Error; err != nil {
		return nil, err
	}

	if cert.Website != nil {
		return &ssl, r.Deploy(cert.ID, cert.WebsiteID)
	}

	return &ssl, nil
}

func (r *certRepo) Renew(id uint) (*acme.Certificate, error) {
	cert, err := r.Get(id)
	if err != nil {
		return nil, err
	}

	client, err := r.getClient(cert)
	if err != nil {
		return nil, err
	}

	if cert.CertURL == "" {
		return nil, errors.New("该证书没有签发成功，无法续签")
	}

	if cert.DNS != nil {
		client.UseDns(acme.DnsType(cert.DNS.Type), cert.DNS.Data)
	} else {
		if cert.Website == nil {
			return nil, errors.New("该证书没有关联网站，无法续签，可以尝试手动签发")
		} else {
			for _, domain := range cert.Domains {
				if strings.Contains(domain, "*") {
					return nil, errors.New("通配符域名无法使用 HTTP 验证")
				}
			}
			conf := fmt.Sprintf("/www/server/vhost/acme/%s.conf", cert.Website.Name)
			client.UseHTTP(conf, cert.Website.Path)
		}
	}

	ssl, err := client.RenewSSL(context.Background(), cert.CertURL, cert.Domains, acme.KeyType(cert.Type))
	if err != nil {
		return nil, err
	}

	cert.CertURL = ssl.URL
	cert.Cert = string(ssl.ChainPEM)
	cert.Key = string(ssl.PrivateKey)
	if err = app.Orm.Save(cert).Error; err != nil {
		return nil, err
	}

	if cert.Website != nil {
		return &ssl, r.Deploy(cert.ID, cert.WebsiteID)
	}

	return &ssl, nil
}

func (r *certRepo) ManualDNS(id uint) ([]acme.DNSRecord, error) {
	cert, err := r.Get(id)
	if err != nil {
		return nil, err
	}

	client, err := r.getClient(cert)
	if err != nil {
		return nil, err
	}

	client.UseManualDns(len(cert.Domains))
	records, err := client.GetDNSRecords(context.Background(), cert.Domains, acme.KeyType(cert.Type))
	if err != nil {
		return nil, err
	}

	// 15 分钟后清理客户端
	r.client = client
	time.AfterFunc(15*time.Minute, func() {
		r.client = nil
	})

	return records, nil
}

func (r *certRepo) Deploy(ID, WebsiteID uint) error {
	cert, err := r.Get(ID)
	if err != nil {
		return err
	}

	if cert.Cert == "" || cert.Key == "" {
		return errors.New("该证书没有签发成功，无法部署")
	}

	website, err := r.websiteRepo.Get(WebsiteID)
	if err != nil {
		return err
	}

	if err = io.Write(fmt.Sprintf("%s/server/vhost/ssl/%s.pem", app.Root, website.Name), cert.Cert, 0644); err != nil {
		return err
	}
	if err = io.Write(fmt.Sprintf("%s/server/vhost/ssl/%s.key", app.Root, website.Name), cert.Key, 0644); err != nil {
		return err
	}
	if err = systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		return err
	}

	return nil
}

func (r *certRepo) getClient(cert *biz.Cert) (*acme.Client, error) {
	var ca string
	var eab *acme.EAB
	switch cert.User.CA {
	case "letsencrypt":
		ca = acme.CALetsEncrypt
	case "buypass":
		ca = acme.CABuypass
	case "zerossl":
		ca = acme.CAZeroSSL
		eab = &acme.EAB{KeyID: cert.User.Kid, MACKey: cert.User.HmacEncoded}
	case "sslcom":
		ca = acme.CASSLcom
		eab = &acme.EAB{KeyID: cert.User.Kid, MACKey: cert.User.HmacEncoded}
	case "google":
		ca = acme.CAGoogle
		eab = &acme.EAB{KeyID: cert.User.Kid, MACKey: cert.User.HmacEncoded}
	}

	return acme.NewPrivateKeyAccount(cert.User.Email, cert.User.PrivateKey, ca, eab)
}
