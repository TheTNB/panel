package acme

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

// GenerateSelfSignedSSL 生成自签名证书
func GenerateSelfSignedSSL(domains []string) ([]byte, []byte, error) {
	rootPrivateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	var ip []net.IP
	isIP := false
	for _, item := range domains {
		ipItem := net.ParseIP(item)
		if len(ipItem) != 0 {
			isIP = true
			ip = append(ip, ipItem)
		}
	}

	rootTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "HaoZi Panel Root CA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(20, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	if isIP {
		rootTemplate.IPAddresses = ip
	} else {
		rootTemplate.DNSNames = domains
	}

	rootCertBytes, _ := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &rootPrivateKey.PublicKey, rootPrivateKey)
	rootCertBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: rootCertBytes,
	}

	interPrivateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	interTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "HaoZi Panel Intermediate CA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(20, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	if isIP {
		interTemplate.IPAddresses = ip
	} else {
		interTemplate.DNSNames = domains
	}

	interCertBytes, _ := x509.CreateCertificate(rand.Reader, &interTemplate, &rootTemplate, &interPrivateKey.PublicKey, rootPrivateKey)
	interCertBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: interCertBytes,
	}

	clientPrivateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	clientTemplate := x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject:      pkix.Name{CommonName: "HaoZi Panel Client"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(20, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	if isIP {
		clientTemplate.IPAddresses = ip
	} else {
		clientTemplate.DNSNames = domains
	}

	clientCertBytes, _ := x509.CreateCertificate(rand.Reader, &clientTemplate, &interTemplate, &clientPrivateKey.PublicKey, interPrivateKey)
	clientCertBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCertBytes,
	}

	pemBytes := []byte{}
	pemBytes = append(pemBytes, pem.EncodeToMemory(clientCertBlock)...)
	pemBytes = append(pemBytes, pem.EncodeToMemory(interCertBlock)...)
	pemBytes = append(pemBytes, pem.EncodeToMemory(rootCertBlock)...)
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientPrivateKey)})

	return pemBytes, keyBytes, nil
}
