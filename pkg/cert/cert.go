package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

func ParseCert(crt string) (x509.Certificate, error) {
	certBlock, _ := pem.Decode([]byte(crt))
	if certBlock == nil {
		return x509.Certificate{}, errors.New("invalid PEM block")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return x509.Certificate{}, err
	}

	return *cert, nil
}

func ParseKey(key string) (crypto.Signer, error) {
	keyBlockDER, _ := pem.Decode([]byte(key))
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

func EncodeCert(cert x509.Certificate) ([]byte, error) {
	pemCert := pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}
	return pem.EncodeToMemory(&pemCert), nil
}

func EncodeKey(key crypto.Signer) ([]byte, error) {
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
