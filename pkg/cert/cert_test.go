package cert

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CertTestSuite struct {
	suite.Suite
}

func TestCertTestSuite(t *testing.T) {
	suite.Run(t, &CertTestSuite{})
}

func (s *CertTestSuite) TestGenerateSelfSigned() {
	pem, key, err := GenerateSelfSigned([]string{"haozi.dev"})
	s.Nil(err)
	s.NotNil(pem)
	s.NotNil(key)
}
