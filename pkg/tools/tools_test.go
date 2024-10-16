package tools

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type HelperTestSuite struct {
	suite.Suite
}

func TestHelperTestSuite(t *testing.T) {
	suite.Run(t, &HelperTestSuite{})
}

func (s *HelperTestSuite) TestGetMonitoringInfo() {
	s.NotNil(CurrentInfo(nil, nil))
}

func (s *HelperTestSuite) TestGetPublicIP() {
	ip, err := GetPublicIP()
	s.Nil(err)
	s.NotEmpty(ip)
}
