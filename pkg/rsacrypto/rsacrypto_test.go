package rsacrypto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RSATestSuite struct {
	suite.Suite
}

func TestRSATestSuite(t *testing.T) {
	suite.Run(t, &RSATestSuite{})
}

func (suite *RSATestSuite) TestRSA() {
	// 生成RSA密钥对
	privateKey, err := GenerateKey()
	suite.NoError(err)
	suite.NotEmpty(privateKey)
	suite.NotEmpty(privateKey.PublicKey)

	// 提取密钥对
	suite.NotEmpty(PrivateKeyToString(privateKey))
	suite.NotEmpty(PublicKeyToString(&privateKey.PublicKey))

	message := []byte("Rat Panel")

	// 加密数据
	ciphertext, err := EncryptData(&privateKey.PublicKey, message)
	suite.NoError(err)
	suite.NotEmpty(ciphertext)

	// 解密数据
	decrypted, err := DecryptData(privateKey, ciphertext)
	suite.NoError(err)
	suite.NotEmpty(decrypted)
}
