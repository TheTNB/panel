package tools

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ToolsTestSuite struct {
	suite.Suite
}

func TestToolsTestSuite(t *testing.T) {
	suite.Run(t, &ToolsTestSuite{})
}

func (s *ToolsTestSuite) TestGetMonitoringInfo() {
	s.NotNil(CurrentInfo(nil, nil))
}

func (s *ToolsTestSuite) TestGetPublicIPv4() {
	ip, err := GetPublicIPv4()
	s.NoError(err)
	s.NotEmpty(ip)
}

func (s *ToolsTestSuite) TestGetPublicIPv6() {
	ip, err := GetPublicIPv6()
	s.Error(err)
	s.Empty(ip)
}

func (s *ToolsTestSuite) TestGetLocalIPv4() {
	ip, err := GetLocalIPv4()
	s.NoError(err)
	s.NotEmpty(ip)
}

func (s *ToolsTestSuite) TestGetLocalIPv6() {
	ip, err := GetLocalIPv6()
	s.Error(err)
	s.Empty(ip)
}

func (s *ToolsTestSuite) TestFormatBytes() {
	s.Equal("1.00 B", FormatBytes(1))
	s.Equal("1.00 KB", FormatBytes(1024))
	s.Equal("1.00 MB", FormatBytes(1024*1024))
	s.Equal("1.00 GB", FormatBytes(1024*1024*1024))
	s.Equal("1.00 TB", FormatBytes(1024*1024*1024*1024))
	s.Equal("1.00 PB", FormatBytes(1024*1024*1024*1024*1024))
	s.Equal("1.00 EB", FormatBytes(1024*1024*1024*1024*1024*1024))
	s.Equal("1.00 ZB", FormatBytes(1024*1024*1024*1024*1024*1024*1024))
	s.Equal("1.00 YB", FormatBytes(1024*1024*1024*1024*1024*1024*1024*1024))
}
