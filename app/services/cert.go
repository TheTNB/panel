// Package services 证书服务
package services

import (
	"errors"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/goravel/framework/facades"
	"panel/pkg/tools"

	requests "panel/app/http/requests/cert"
	"panel/app/models"
	"panel/pkg/acme"
)

type Cert interface {
	GetByID(ID uint) (models.Cert, error)
	UserAdd(request requests.UserAdd) error
	UserDelete(ID uint) error
	DNSAdd(request requests.DNSAdd) error
	DNSDelete(ID uint) error
	CertAdd(request requests.CertAdd) error
	CertDelete(ID uint) error
	ObtainAuto(ID uint) (certificate.Resource, error)
	ObtainManual(ID uint) (certificate.Resource, error)
	ManualDNS(ID uint) (map[string]acme.Resolve, error)
	Renew(ID uint) (certificate.Resource, error)
}

type CertImpl struct {
}

func NewCertImpl() *CertImpl {
	return &CertImpl{}
}

// GetByID 根据 ID 获取证书
func (s *CertImpl) GetByID(ID uint) (models.Cert, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("User").With("DNS").With("Website").Where("id = ?", ID).First(&cert)
	return cert, err
}

// UserAdd 添加用户
func (s *CertImpl) UserAdd(request requests.UserAdd) error {
	var user models.CertUser
	user.CA = request.CA
	user.Email = request.Email
	user.Kid = &request.Kid
	user.HmacEncoded = &request.HmacEncoded
	user.KeyType = request.KeyType

	var err error
	var client *acme.Client
	switch user.CA {
	case "letsencrypt":
		client, err = acme.NewRegisterClient(user.Email, acme.CALetEncrypt, certcrypto.KeyType(user.KeyType))
	case "buypass":
		client, err = acme.NewRegisterClient(user.Email, acme.CABuypass, certcrypto.KeyType(user.KeyType))
	case "zerossl":
		client, err = acme.NewRegisterWithExternalAccountBindingClient(user.Email, *user.Kid, *user.HmacEncoded, acme.CAZeroSSL, certcrypto.KeyType(user.KeyType))
	case "sslcom":
		client, err = acme.NewRegisterWithExternalAccountBindingClient(user.Email, *user.Kid, *user.HmacEncoded, acme.CASSLcom, certcrypto.KeyType(user.KeyType))
	case "google":
		client, err = acme.NewRegisterWithExternalAccountBindingClient(user.Email, *user.Kid, *user.HmacEncoded, acme.CAGoogle, certcrypto.KeyType(user.KeyType))
	default:
		return errors.New("CA 提供商不支持")
	}

	if err != nil {
		return errors.New("向 CA 注册账号失败，请检查参数是否正确")
	}

	privateKey, err := acme.GetPrivateKey(client.User.GetPrivateKey(), acme.KeyType(user.KeyType))
	if err != nil {
		return errors.New("获取私钥失败")
	}
	user.PrivateKey = string(privateKey)

	return facades.Orm().Query().Create(&user)
}

// UserDelete 删除用户
func (s *CertImpl) UserDelete(ID uint) error {
	var user models.CertUser
	err := facades.Orm().Query().With("Certs").Where("id = ?", ID).First(&user)
	if err != nil {
		return err
	}

	if user.Certs != nil {
		return errors.New("该用户下存在证书，无法删除")
	}

	_, err = facades.Orm().Query().Delete(&models.CertUser{}, ID)
	return err
}

// DNSAdd 添加 DNS
func (s *CertImpl) DNSAdd(request requests.DNSAdd) error {
	var dns models.CertDNS
	dns.Type = request.Type
	dns.Data = request.Data

	return facades.Orm().Query().Create(&dns)
}

// DNSDelete 删除 DNS
func (s *CertImpl) DNSDelete(ID uint) error {
	var dns models.CertDNS
	err := facades.Orm().Query().With("Certs").Where("id = ?", ID).First(&dns)
	if err != nil {
		return err
	}

	if dns.Certs != nil {
		return errors.New("该 DNS 接口下存在证书，无法删除")
	}

	_, err = facades.Orm().Query().Delete(&models.CertDNS{}, ID)
	return err
}

// CertAdd 添加证书
func (s *CertImpl) CertAdd(request requests.CertAdd) error {
	var cert models.Cert
	cert.Type = request.Type
	cert.Domains = request.Domains
	cert.UserID = request.UserID

	if request.DNSID != nil {
		cert.DNSID = request.DNSID
		// TODO 生成计划任务
	}

	return facades.Orm().Query().Create(&cert)
}

// CertDelete 删除证书
func (s *CertImpl) CertDelete(ID uint) error {
	var cert models.Cert
	err := facades.Orm().Query().Where("id = ?", ID).First(&cert)
	if err != nil {
		return err
	}

	if cert.CronID != nil {
		// TODO 删除计划任务
	}

	_, err = facades.Orm().Query().Delete(&models.Cert{}, ID)
	return err
}

// ObtainAuto 自动签发证书
func (s *CertImpl) ObtainAuto(ID uint) (certificate.Resource, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("Website").With("User").With("DNS").Where("id = ?", ID).First(&cert)
	if err != nil {
		return certificate.Resource{}, err
	}

	var ca string
	switch cert.User.CA {
	case "letsencrypt":
		ca = acme.CALetEncrypt
	case "buypass":
		ca = acme.CABuypass
	case "zerossl":
		ca = acme.CAZeroSSL
	case "sslcom":
		ca = acme.CASSLcom
	case "google":
		ca = acme.CAGoogle
	}

	client, err := acme.NewPrivateKeyClient(cert.User.Email, cert.User.PrivateKey, ca, certcrypto.KeyType(cert.User.KeyType))
	if err != nil {
		return certificate.Resource{}, err
	}

	if cert.DNS != nil {
		err = client.UseDns(acme.DnsType(cert.DNS.Type), cert.DNS.Data)
	} else {
		if cert.Website == nil {
			return certificate.Resource{}, errors.New("该证书没有关联网站，无法自动签发")
		} else {
			err = client.UseHTTP(cert.Website.Path)
		}
	}
	if err != nil {
		return certificate.Resource{}, err
	}

	ssl, err := client.ObtainSSL(cert.Domains)
	if err != nil {
		return certificate.Resource{}, err
	}

	cert.CertURL = &ssl.CertURL
	cert.Cert = string(ssl.Certificate)
	cert.Key = string(ssl.PrivateKey)
	err = facades.Orm().Query().Save(&cert)
	if err != nil {
		return certificate.Resource{}, err
	}

	if cert.Website != nil {
		tools.Write("/www/server/vhost/ssl/"+cert.Website.Name+".pem", string(ssl.Certificate), 0644)
		tools.Write("/www/server/vhost/ssl/"+cert.Website.Name+".key", string(ssl.PrivateKey), 0644)
		tools.Exec("systemctl reload openresty")
	}

	return ssl, nil
}

// ObtainManual 手动签发证书
func (s *CertImpl) ObtainManual(ID uint) (certificate.Resource, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("User").Where("id = ?", ID).First(&cert)
	if err != nil {
		return certificate.Resource{}, err
	}

	var ca string
	switch cert.User.CA {
	case "letsencrypt":
		ca = acme.CALetEncrypt
	case "buypass":
		ca = acme.CABuypass
	case "zerossl":
		ca = acme.CAZeroSSL
	case "sslcom":
		ca = acme.CASSLcom
	case "google":
		ca = acme.CAGoogle
	}

	client, err := acme.NewPrivateKeyClient(cert.User.Email, cert.User.PrivateKey, ca, certcrypto.KeyType(cert.User.KeyType))
	if err != nil {
		return certificate.Resource{}, err
	}

	err = client.UseManualDns()
	if err != nil {
		return certificate.Resource{}, err
	}

	ssl, err := client.ObtainSSL(cert.Domains)
	if err != nil {
		return certificate.Resource{}, err
	}

	cert.CertURL = &ssl.CertURL
	cert.Cert = string(ssl.Certificate)
	cert.Key = string(ssl.PrivateKey)
	err = facades.Orm().Query().Save(&cert)
	if err != nil {
		return certificate.Resource{}, err
	}

	if cert.Website != nil {
		tools.Write("/www/server/vhost/ssl/"+cert.Website.Name+".pem", string(ssl.Certificate), 0644)
		tools.Write("/www/server/vhost/ssl/"+cert.Website.Name+".key", string(ssl.PrivateKey), 0644)
		tools.Exec("systemctl reload openresty")
	}

	return ssl, nil
}

// ManualDNS 获取手动 DNS 解析信息
func (s *CertImpl) ManualDNS(ID uint) (map[string]acme.Resolve, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("User").Where("id = ?", ID).First(&cert)
	if err != nil {
		return nil, err
	}

	var ca string
	switch cert.User.CA {
	case "letsencrypt":
		ca = acme.CALetEncrypt
	case "buypass":
		ca = acme.CABuypass
	case "zerossl":
		ca = acme.CAZeroSSL
	case "sslcom":
		ca = acme.CASSLcom
	case "google":
		ca = acme.CAGoogle
	}

	client, err := acme.NewPrivateKeyClient(cert.User.Email, cert.User.PrivateKey, ca, certcrypto.KeyType(cert.User.KeyType))
	if err != nil {
		return nil, err
	}

	err = client.UseManualDns()
	if err != nil {
		return nil, err
	}

	return client.GetDNSResolve(cert.Domains)
}

// Renew 续签证书
func (s *CertImpl) Renew(ID uint) (certificate.Resource, error) {
	var cert models.Cert
	err := facades.Orm().Query().With("User").With("DNS").Where("id = ?", ID).First(&cert)
	if err != nil {
		return certificate.Resource{}, err
	}

	var ca string
	switch cert.User.CA {
	case "letsencrypt":
		ca = acme.CALetEncrypt
	case "buypass":
		ca = acme.CABuypass
	case "zerossl":
		ca = acme.CAZeroSSL
	case "sslcom":
		ca = acme.CASSLcom
	case "google":
		ca = acme.CAGoogle
	}

	client, err := acme.NewPrivateKeyClient(cert.User.Email, cert.User.PrivateKey, ca, certcrypto.KeyType(cert.User.KeyType))
	if err != nil {
		return certificate.Resource{}, err
	}

	if cert.CertURL == nil {
		return certificate.Resource{}, errors.New("该证书没有签发成功，无法续签")
	}

	if cert.DNS != nil {
		err = client.UseDns(acme.DnsType(cert.DNS.Type), cert.DNS.Data)
	} else {
		if cert.Website == nil {
			return certificate.Resource{}, errors.New("该证书没有关联网站，无法续签，可以尝试手动签发")
		} else {
			err = client.UseHTTP(cert.Website.Path)
		}
	}
	if err != nil {
		return certificate.Resource{}, err
	}

	ssl, err := client.RenewSSL(*cert.CertURL)
	if err != nil {
		return certificate.Resource{}, err
	}

	cert.CertURL = &ssl.CertURL
	cert.Cert = string(ssl.Certificate)
	cert.Key = string(ssl.PrivateKey)
	err = facades.Orm().Query().Save(&cert)
	if err != nil {
		return certificate.Resource{}, err
	}

	if cert.Website != nil {
		tools.Write("/www/server/vhost/ssl/"+cert.Website.Name+".pem", string(ssl.Certificate), 0644)
		tools.Write("/www/server/vhost/ssl/"+cert.Website.Name+".key", string(ssl.PrivateKey), 0644)
		tools.Exec("systemctl reload openresty")
	}

	return ssl, nil
}
