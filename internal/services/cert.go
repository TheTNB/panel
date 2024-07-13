// Package services 证书服务
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/goravel/framework/facades"

	requests "github.com/TheTNB/panel/v2/app/http/requests/cert"
	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/pkg/acme"
	"github.com/TheTNB/panel/v2/pkg/cert"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type CertImpl struct {
	client *acme.Client
}

func NewCertImpl() *CertImpl {
	return &CertImpl{}
}

// UserStore 添加用户
func (s *CertImpl) UserStore(request requests.UserStore) error {
	var user models.CertUser
	user.CA = request.CA
	user.Email = request.Email
	user.Kid = request.Kid
	user.HmacEncoded = request.HmacEncoded
	user.KeyType = request.KeyType

	var err error
	var client *acme.Client
	switch user.CA {
	case "letsencrypt":
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CALetsEncrypt, nil, acme.KeyType(user.KeyType))
	case "buypass":
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CABuypass, nil, acme.KeyType(user.KeyType))
	case "zerossl":
		eab, eabErr := s.getZeroSSLEAB(user.Email)
		if eabErr != nil {
			return eabErr
		}
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CAZeroSSL, eab, acme.KeyType(user.KeyType))
	case "sslcom":
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CASSLcom, &acme.EAB{KeyID: user.Kid, MACKey: user.HmacEncoded}, acme.KeyType(user.KeyType))
	case "google":
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CAGoogle, &acme.EAB{KeyID: user.Kid, MACKey: user.HmacEncoded}, acme.KeyType(user.KeyType))
	default:
		return errors.New("CA 提供商不支持")
	}

	if err != nil {
		return errors.New("向 CA 注册账号失败，请检查参数是否正确")
	}

	privateKey, err := cert.EncodeKey(client.Account.PrivateKey)
	if err != nil {
		return errors.New("获取私钥失败")
	}
	user.PrivateKey = string(privateKey)

	return facades.Orm().Query().Create(&user)
}

// UserUpdate 更新用户
func (s *CertImpl) UserUpdate(request requests.UserUpdate) error {
	var user models.CertUser
	err := facades.Orm().Query().Where("id = ?", request.ID).First(&user)
	if err != nil {
		return err
	}

	user.CA = request.CA
	user.Email = request.Email
	user.Kid = request.Kid
	user.HmacEncoded = request.HmacEncoded
	user.KeyType = request.KeyType

	var client *acme.Client
	switch user.CA {
	case "letsencrypt":
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CALetsEncrypt, nil, acme.KeyType(user.KeyType))
	case "buypass":
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CABuypass, nil, acme.KeyType(user.KeyType))
	case "zerossl":
		eab, eabErr := s.getZeroSSLEAB(user.Email)
		if eabErr != nil {
			return eabErr
		}
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CAZeroSSL, eab, acme.KeyType(user.KeyType))
	case "sslcom":
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CASSLcom, &acme.EAB{KeyID: user.Kid, MACKey: user.HmacEncoded}, acme.KeyType(user.KeyType))
	case "google":
		client, err = acme.NewRegisterAccount(context.Background(), user.Email, acme.CAGoogle, &acme.EAB{KeyID: user.Kid, MACKey: user.HmacEncoded}, acme.KeyType(user.KeyType))
	default:
		return errors.New("CA 提供商不支持")
	}

	if err != nil {
		return errors.New("向 CA 注册账号失败，请检查参数是否正确")
	}

	privateKey, err := cert.EncodeKey(client.Account.PrivateKey)
	if err != nil {
		return errors.New("获取私钥失败")
	}
	user.PrivateKey = string(privateKey)

	return facades.Orm().Query().Save(&user)
}

// getZeroSSLEAB 获取 ZeroSSL EAB
func (s *CertImpl) getZeroSSLEAB(email string) (*acme.EAB, error) {
	type data struct {
		Success    bool   `json:"success"`
		EabKid     string `json:"eab_kid"`
		EabHmacKey string `json:"eab_hmac_key"`
	}
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	resp, err := client.R().SetFormData(map[string]string{
		"email": email,
	}).SetResult(&data{}).Post("https://api.zerossl.com/acme/eab-credentials-email")
	if err != nil || !resp.IsSuccess() {
		return &acme.EAB{}, errors.New("获取ZeroSSL EAB失败")
	}
	eab := resp.Result().(*data)
	if !eab.Success {
		return &acme.EAB{}, errors.New("获取ZeroSSL EAB失败")
	}

	return &acme.EAB{KeyID: eab.EabKid, MACKey: eab.EabHmacKey}, nil
}

// UserShow 根据 ID 获取用户
func (s *CertImpl) UserShow(ID uint) (models.CertUser, error) {
	var user models.CertUser
	err := facades.Orm().Query().With("Certs").Where("id = ?", ID).First(&user)

	return user, err
}

// UserDestroy 删除用户
func (s *CertImpl) UserDestroy(ID uint) error {
	var cert models.Cert
	err := facades.Orm().Query().Where("user_id = ?", ID).First(&cert)
	if err != nil {
		return err
	}

	if cert.ID != 0 {
		return errors.New("该用户下存在证书，无法删除")
	}

	_, err = facades.Orm().Query().Delete(&models.CertUser{}, ID)
	return err
}

// DNSStore 添加 DNS
func (s *CertImpl) DNSStore(request requests.DNSStore) error {
	var dns models.CertDNS
	dns.Type = request.Type
	dns.Name = request.Name
	dns.Data = request.Data

	return facades.Orm().Query().Create(&dns)
}

// DNSUpdate 更新 DNS
func (s *CertImpl) DNSUpdate(request requests.DNSUpdate) error {
	var dns models.CertDNS
	err := facades.Orm().Query().Where("id = ?", request.ID).First(&dns)
	if err != nil {
		return err
	}

	dns.Type = request.Type
	dns.Name = request.Name
	dns.Data = request.Data

	return facades.Orm().Query().Save(&dns)
}

// DNSShow 根据 ID 获取 DNS
func (s *CertImpl) DNSShow(ID uint) (models.CertDNS, error) {
	var dns models.CertDNS
	err := facades.Orm().Query().With("Certs").Where("id = ?", ID).First(&dns)

	return dns, err
}

// DNSDestroy 删除 DNS
func (s *CertImpl) DNSDestroy(ID uint) error {
	var cert models.Cert
	err := facades.Orm().Query().Where("dns_id = ?", ID).First(&cert)
	if err != nil {
		return err
	}

	if cert.ID != 0 {
		return errors.New("该 DNS 接口下存在证书，无法删除")
	}

	_, err = facades.Orm().Query().Delete(&models.CertDNS{}, ID)
	return err
}

// CertStore 添加证书
func (s *CertImpl) CertStore(request requests.CertStore) error {
	var cert models.Cert
	cert.Type = request.Type
	cert.Domains = request.Domains
	cert.AutoRenew = request.AutoRenew
	cert.UserID = request.UserID
	cert.DNSID = request.DNSID
	cert.WebsiteID = request.WebsiteID

	return facades.Orm().Query().Create(&cert)
}

// CertUpdate 更新证书
func (s *CertImpl) CertUpdate(request requests.CertUpdate) error {
	var cert models.Cert
	err := facades.Orm().Query().Where("id = ?", request.ID).First(&cert)
	if err != nil {
		return err
	}

	cert.Type = request.Type
	cert.Domains = request.Domains
	cert.AutoRenew = request.AutoRenew
	cert.UserID = request.UserID
	cert.DNSID = request.DNSID
	cert.WebsiteID = request.WebsiteID

	return facades.Orm().Query().Save(&cert)
}

// CertShow 根据 ID 获取证书
func (s *CertImpl) CertShow(ID uint) (models.Cert, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("User").With("DNS").With("Website").Where("id = ?", ID).First(&cert)

	return cert, err
}

// CertDestroy 删除证书
func (s *CertImpl) CertDestroy(ID uint) error {
	var cert models.Cert
	err := facades.Orm().Query().Where("id = ?", ID).First(&cert)
	if err != nil {
		return err
	}

	_, err = facades.Orm().Query().Delete(&models.Cert{}, ID)
	return err
}

// ObtainAuto 自动签发证书
func (s *CertImpl) ObtainAuto(ID uint) (acme.Certificate, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("Website").With("User").With("DNS").Where("id = ?", ID).First(&cert)
	if err != nil {
		return acme.Certificate{}, err
	}

	client, err := s.getClient(cert)
	if err != nil {
		return acme.Certificate{}, err
	}

	if cert.DNS != nil {
		client.UseDns(acme.DnsType(cert.DNS.Type), cert.DNS.Data)
	} else {
		if cert.Website == nil {
			return acme.Certificate{}, errors.New("该证书没有关联网站，无法自动签发")
		} else {
			for _, domain := range cert.Domains {
				if strings.Contains(domain, "*") {
					return acme.Certificate{}, errors.New("通配符域名无法使用 HTTP 验证")
				}
			}
			conf := fmt.Sprintf("/www/server/vhost/acme/%s.conf", cert.Website.Name)
			client.UseHTTP(conf, cert.Website.Path)
		}
	}

	ssl, err := client.ObtainSSL(context.Background(), cert.Domains, acme.KeyType(cert.Type))
	if err != nil {
		return acme.Certificate{}, err
	}

	cert.CertURL = ssl.URL
	cert.Cert = string(ssl.ChainPEM)
	cert.Key = string(ssl.PrivateKey)
	err = facades.Orm().Query().Save(&cert)
	if err != nil {
		return acme.Certificate{}, err
	}

	if cert.Website != nil {
		if err = io.Write("/www/server/vhost/ssl/"+cert.Website.Name+".pem", cert.Cert, 0644); err != nil {
			return acme.Certificate{}, err
		}
		if err = io.Write("/www/server/vhost/ssl/"+cert.Website.Name+".key", cert.Key, 0644); err != nil {
			return acme.Certificate{}, err
		}
		if err = systemctl.Reload("openresty"); err != nil {
			_, err = shell.Execf("openresty -t")
			return acme.Certificate{}, err
		}
	}

	return ssl, nil
}

// ObtainManual 手动签发证书
func (s *CertImpl) ObtainManual(ID uint) (acme.Certificate, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("User").Where("id = ?", ID).First(&cert)
	if err != nil {
		return acme.Certificate{}, err
	}

	if s.client == nil {
		return acme.Certificate{}, errors.New("请重新获取 DNS 解析记录")
	}

	ssl, err := s.client.ObtainSSLManual()
	if err != nil {
		return acme.Certificate{}, err
	}

	cert.CertURL = ssl.URL
	cert.Cert = string(ssl.ChainPEM)
	cert.Key = string(ssl.PrivateKey)
	err = facades.Orm().Query().Save(&cert)
	if err != nil {
		return acme.Certificate{}, err
	}

	if cert.Website != nil {
		if err = io.Write("/www/server/vhost/ssl/"+cert.Website.Name+".pem", cert.Cert, 0644); err != nil {
			return acme.Certificate{}, err
		}
		if err = io.Write("/www/server/vhost/ssl/"+cert.Website.Name+".key", cert.Key, 0644); err != nil {
			return acme.Certificate{}, err
		}
		if err = systemctl.Reload("openresty"); err != nil {
			_, err = shell.Execf("openresty -t")
			return acme.Certificate{}, err
		}
	}

	return ssl, nil
}

// ManualDNS 获取手动 DNS 解析信息
func (s *CertImpl) ManualDNS(ID uint) ([]acme.DNSRecord, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("User").Where("id = ?", ID).First(&cert)
	if err != nil {
		return nil, err
	}

	client, err := s.getClient(cert)
	if err != nil {
		return nil, err
	}

	client.UseManualDns(len(cert.Domains))
	records, err := client.GetDNSRecords(context.Background(), cert.Domains, acme.KeyType(cert.Type))

	// 15 分钟后清理客户端
	s.client = client
	time.AfterFunc(15*time.Minute, func() {
		s.client = nil
	})

	return records, err
}

// Renew 续签证书
func (s *CertImpl) Renew(ID uint) (acme.Certificate, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("Website").With("User").With("DNS").Where("id = ?", ID).First(&cert)
	if err != nil {
		return acme.Certificate{}, err
	}

	client, err := s.getClient(cert)
	if err != nil {
		return acme.Certificate{}, err
	}

	if cert.CertURL == "" {
		return acme.Certificate{}, errors.New("该证书没有签发成功，无法续签")
	}

	if cert.DNS != nil {
		client.UseDns(acme.DnsType(cert.DNS.Type), cert.DNS.Data)
	} else {
		if cert.Website == nil {
			return acme.Certificate{}, errors.New("该证书没有关联网站，无法续签，可以尝试手动签发")
		} else {
			for _, domain := range cert.Domains {
				if strings.Contains(domain, "*") {
					return acme.Certificate{}, errors.New("通配符域名无法使用 HTTP 验证")
				}
			}
			conf := fmt.Sprintf("/www/server/vhost/acme/%s.conf", cert.Website.Name)
			client.UseHTTP(conf, cert.Website.Path)
		}
	}

	ssl, err := client.RenewSSL(context.Background(), cert.CertURL, cert.Domains, acme.KeyType(cert.Type))
	if err != nil {
		return acme.Certificate{}, err
	}

	cert.CertURL = ssl.URL
	cert.Cert = string(ssl.ChainPEM)
	cert.Key = string(ssl.PrivateKey)
	err = facades.Orm().Query().Save(&cert)
	if err != nil {
		return acme.Certificate{}, err
	}

	if cert.Website != nil {
		if err = io.Write("/www/server/vhost/ssl/"+cert.Website.Name+".pem", cert.Cert, 0644); err != nil {
			return acme.Certificate{}, err
		}
		if err = io.Write("/www/server/vhost/ssl/"+cert.Website.Name+".key", cert.Key, 0644); err != nil {
			return acme.Certificate{}, err
		}
		if err = systemctl.Reload("openresty"); err != nil {
			_, err = shell.Execf("openresty -t")
			return acme.Certificate{}, err
		}
	}

	return ssl, nil
}

// Deploy 部署证书
func (s *CertImpl) Deploy(ID, WebsiteID uint) error {
	var cert models.Cert
	err := facades.Orm().Query().Where("id = ?", ID).First(&cert)
	if err != nil {
		return err
	}

	if cert.Cert == "" || cert.Key == "" {
		return errors.New("该证书没有签发成功，无法部署")
	}

	website := models.Website{}
	err = facades.Orm().Query().Where("id = ?", WebsiteID).First(&website)
	if err != nil {
		return err
	}

	if err = io.Write("/www/server/vhost/ssl/"+website.Name+".pem", cert.Cert, 0644); err != nil {
		return err
	}
	if err = io.Write("/www/server/vhost/ssl/"+website.Name+".key", cert.Key, 0644); err != nil {
		return err
	}
	if err = systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		return err
	}

	return nil
}

func (s *CertImpl) getClient(cert models.Cert) (*acme.Client, error) {
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
