package acme

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net/http"

	"github.com/mholt/acmez/v2"
	"github.com/mholt/acmez/v2/acme"
	"go.uber.org/zap"

	"github.com/TheTNB/panel/v2/pkg/cert"
)

const (
	CALetsEncryptStaging = "https://acme-staging-v02.api.letsencrypt.org/directory"
	CALetsEncrypt        = "https://acme-v02.api.letsencrypt.org/directory"
	CAZeroSSL            = "https://acme.zerossl.com/v2/DV90"
	CAGoogle             = "https://dv.acme-v02.api.pki.goog/directory"
	CABuypass            = "https://api.buypass.com/acme/directory"
	CASSLcom             = "https://acme.ssl.com/sslcom-dv-rsa"
)

type KeyType string

const (
	KeyEC256   = KeyType("P256")
	KeyEC384   = KeyType("P384")
	KeyRSA2048 = KeyType("2048")
	KeyRSA3072 = KeyType("3072")
	KeyRSA4096 = KeyType("4096")
)

type EAB = acme.EAB

func NewRegisterAccount(ctx context.Context, email, CA string, eab *EAB, keyType KeyType) (*Client, error) {
	client, err := getClient(CA)
	if err != nil {
		return nil, err
	}

	accountPrivateKey, err := generatePrivateKey(keyType)
	if err != nil {
		return nil, err
	}
	account := acme.Account{
		Contact:              []string{"mailto:" + email},
		TermsOfServiceAgreed: true,
		PrivateKey:           accountPrivateKey,
	}
	if eab != nil {
		err = account.SetExternalAccountBinding(ctx, client.Client, *eab)
		if err != nil {
			return nil, err
		}
	}

	account, err = client.NewAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	return &Client{Account: account, zClient: client}, nil
}

func NewPrivateKeyAccount(email string, privateKey string, CA string, eab *EAB) (*Client, error) {
	client, err := getClient(CA)
	if err != nil {
		return nil, err
	}

	key, err := cert.ParseKey(privateKey)
	if err != nil {
		return nil, err
	}

	account := acme.Account{
		Contact:              []string{"mailto:" + email},
		TermsOfServiceAgreed: true,
		PrivateKey:           key,
	}
	if eab != nil {
		err = account.SetExternalAccountBinding(context.Background(), client.Client, *eab)
		if err != nil {
			return nil, err
		}
	}

	account, err = client.GetAccount(context.Background(), account)
	if err != nil {
		return nil, err
	}

	return &Client{Account: account, zClient: client}, nil
}

func generatePrivateKey(keyType KeyType) (crypto.Signer, error) {
	switch keyType {
	case KeyEC256:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case KeyEC384:
		return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case KeyRSA2048:
		return rsa.GenerateKey(rand.Reader, 2048)
	case KeyRSA3072:
		return rsa.GenerateKey(rand.Reader, 3072)
	case KeyRSA4096:
		return rsa.GenerateKey(rand.Reader, 4096)
	}

	return nil, errors.New("未知的密钥类型")
}

func getClient(CA string) (acmez.Client, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return acmez.Client{}, err
	}

	client := acmez.Client{
		Client: &acme.Client{
			Directory:  CA,
			HTTPClient: http.DefaultClient,
			Logger:     logger,
		},
	}

	return client, nil
}
