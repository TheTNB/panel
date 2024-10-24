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

func (s *HelperTestSuite) TestGetPublicIPv4() {
	ip, err := GetPublicIPv4()
	s.Nil(err)
	s.NotEmpty(ip)
}

func (s *HelperTestSuite) TestGetPublicIPv6() {
	ip, err := GetPublicIPv6()
	s.Nil(err)
	s.NotEmpty(ip)
}

func (s *HelperTestSuite) TestGetLocalIPv4() {
	ip, err := GetLocalIPv4()
	s.Nil(err)
	s.NotEmpty(ip)
}

func (s *HelperTestSuite) TestGetLocalIPv6() {
	ip, err := GetLocalIPv6()
	s.Nil(err)
	s.NotEmpty(ip)
}
