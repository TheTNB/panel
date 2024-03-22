package acme

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mholt/acmez"
	"github.com/mholt/acmez/acme"
	"go.uber.org/zap"
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

	key, err := parsePrivateKey([]byte(privateKey))
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

func parsePrivateKey(key []byte) (crypto.Signer, error) {
	keyBlockDER, _ := pem.Decode(key)
	if keyBlockDER == nil {
		return nil, errors.New("invalid PEM block")
	}

	if keyBlockDER.Type != "PRIVATE KEY" && !strings.HasSuffix(keyBlockDER.Type, " PRIVATE KEY") {
		return nil, fmt.Errorf("unknown PEM header %q", keyBlockDER.Type)
	}

	if parse, err := x509.ParsePKCS1PrivateKey(keyBlockDER.Bytes); err == nil {
		return parse, nil
	}

	if parse, err := x509.ParsePKCS8PrivateKey(keyBlockDER.Bytes); err == nil {
		switch parse.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			return parse.(crypto.Signer), nil
		default:
			return nil, fmt.Errorf("found unknown private key type in PKCS#8 wrapping: %T", key)
		}
	}

	if parse, err := x509.ParseECPrivateKey(keyBlockDER.Bytes); err == nil {
		return parse, nil
	}

	return nil, errors.New("解析私钥失败")
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

func EncodePrivateKey(key crypto.Signer) ([]byte, error) {
	var pemType string
	var keyBytes []byte
	switch key := key.(type) {
	case *ecdsa.PrivateKey:
		var err error
		pemType = "EC"
		keyBytes, err = x509.MarshalECPrivateKey(key)
		if err != nil {
			return nil, err
		}
	case *rsa.PrivateKey:
		pemType = "RSA"
		keyBytes = x509.MarshalPKCS1PrivateKey(key)
	case ed25519.PrivateKey:
		var err error
		pemType = "ED25519"
		keyBytes, err = x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("未知的密钥类型 %T", key)
	}
	pemKey := pem.Block{Type: pemType + " PRIVATE KEY", Bytes: keyBytes}
	return pem.EncodeToMemory(&pemKey), nil
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
