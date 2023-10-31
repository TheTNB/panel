package acme

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

const (
	CALetEncrypt = "https://acme-v02.api.letsencrypt.org/directory"
	CAZeroSSL    = "https://acme.zerossl.com/v2/DV90"
	CAGoogle     = "https://dv.acme-v02.api.pki.goog/directory"
	CABuypass    = "https://api.buypass.com/acme/directory"
	CASSLcom     = "https://acme.ssl.com/sslcom-dv-rsa"
)

const (
	KeyEC256   = certcrypto.EC256
	KeyEC384   = certcrypto.EC384
	KeyRSA2048 = certcrypto.RSA2048
	KeyRSA3072 = certcrypto.RSA3072
	KeyRSA4096 = certcrypto.RSA4096
)

type domainError struct {
	Domain string
	Error  error
}

type User struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

func GetPrivateKey(priKey crypto.PrivateKey) []byte {
	rsaKey := priKey.(*rsa.PrivateKey)
	marshal := x509.MarshalPKCS1PrivateKey(rsaKey)
	block := &pem.Block{
		Type:  "privateKey",
		Bytes: marshal,
	}
	return pem.EncodeToMemory(block)
}

func NewRegisterClient(email string, CA string, keyType certcrypto.KeyType) (*Client, error) {
	privateKey, err := certcrypto.GeneratePrivateKey(keyType)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email: email,
		Key:   privateKey,
	}
	config := lego.NewConfig(user)
	config.CADirURL = CA
	config.Certificate.KeyType = keyType
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	user.Registration = reg

	acmeClient := &Client{
		User:   user,
		Client: client,
		Config: config,
	}

	return acmeClient, nil
}

func NewRegisterWithExternalAccountBindingClient(email, kid, hmac, CA string, keyType certcrypto.KeyType) (*Client, error) {
	privateKey, err := certcrypto.GeneratePrivateKey(keyType)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email: email,
		Key:   privateKey,
	}
	config := lego.NewConfig(user)
	config.CADirURL = CA
	config.Certificate.KeyType = keyType
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	reg, err := client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{TermsOfServiceAgreed: true, Kid: kid, HmacEncoded: hmac})
	if err != nil {
		return nil, err
	}
	user.Registration = reg

	acmeClient := &Client{
		User:   user,
		Client: client,
		Config: config,
	}

	return acmeClient, nil
}

func NewPrivateKeyClient(email string, privateKey string, CA string, keyType certcrypto.KeyType) (*Client, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("无法解析私钥")
	}

	key, err := certcrypto.ParsePEMPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email: email,
		Key:   key,
	}
	config := lego.NewConfig(user)
	config.CADirURL = CA
	config.Certificate.KeyType = keyType
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	reg, err := client.Registration.ResolveAccountByKey()
	if err != nil {
		return nil, err
	}
	user.Registration = reg

	acmeClient := &Client{
		User:   user,
		Client: client,
		Config: config,
	}

	return acmeClient, nil
}
