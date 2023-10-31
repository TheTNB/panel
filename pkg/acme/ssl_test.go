package acme

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SSLTestSuite struct {
	suite.Suite
}

func TestSSLTestSuite(t *testing.T) {
	suite.Run(t, &SSLTestSuite{})
}

func (s *SSLTestSuite) TestGenerateSelfSignedSSL() {
	pem, key, err := GenerateSelfSignedSSL([]string{"haozi.dev"})
	s.Nil(err)
	s.NotNil(pem)
	s.NotNil(key)
}
