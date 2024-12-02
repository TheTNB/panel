package rsacrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

const (
	keySize = 2048 // RSA key size in bits
)

// GenerateKey 生成RSA密钥对
func GenerateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, keySize)
}

// EncryptData 加密数据
func EncryptData(publicKey *rsa.PublicKey, data []byte) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		publicKey,
		data,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %v", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptData 解密数据
func DecryptData(privateKey *rsa.PrivateKey, ciphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	plaintext, err := rsa.DecryptOAEP(
		sha512.New(),
		rand.Reader,
		privateKey,
		data,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}

	return plaintext, nil
}

// PrivateKeyToString 将RSA私钥转换为PEM格式的字符串
func PrivateKeyToString(privateKey *rsa.PrivateKey) (string, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	return string(privateKeyPEM), nil
}

// PublicKeyToString 将RSA公钥转换为PEM格式的字符串
func PublicKeyToString(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %v", err)
	}

	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	return string(publicKeyPEM), nil
}
