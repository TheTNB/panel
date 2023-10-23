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
	s.NotNil(GetMonitoringInfo())
}

func (s *HelperTestSuite) TestVersionCompare() {
	// 测试相等情况
	s.True(VersionCompare("1.0.0", "1.0.0", "=="))
	s.True(VersionCompare("1.0.0", "1.0.0", ">="))
	s.True(VersionCompare("1.0.0", "1.0.0", "<="))
	s.False(VersionCompare("1.0.0", "1.0.0", ">"))
	s.False(VersionCompare("1.0.0", "1.0.0", "<"))
	s.False(VersionCompare("1.0.0", "1.0.0", "!="))

	// 测试1.0.0小于1.0.1
	s.True(VersionCompare("1.0.0", "1.0.1", "<"))
	s.True(VersionCompare("1.0.0", "1.0.1", "<="))
	s.True(VersionCompare("1.0.0", "1.0.1", "!="))
	s.False(VersionCompare("1.0.0", "1.0.1", "=="))
	s.False(VersionCompare("1.0.0", "1.0.1", ">="))
	s.False(VersionCompare("1.0.0", "1.0.1", ">"))

	// 测试1.0.1大于1.0.0
	s.True(VersionCompare("1.0.1", "1.0.0", ">"))
	s.True(VersionCompare("1.0.1", "1.0.0", ">="))
	s.True(VersionCompare("1.0.1", "1.0.0", "!="))
	s.False(VersionCompare("1.0.1", "1.0.0", "=="))
	s.False(VersionCompare("1.0.1", "1.0.0", "<="))
	s.False(VersionCompare("1.0.1", "1.0.0", "<"))

	// 测试带有 'v' 前缀的版本号
	s.True(VersionCompare("v1.0.0", "1.0.0", "=="))
	s.True(VersionCompare("1.0.0", "v1.0.0", "=="))
	s.True(VersionCompare("v1.0.0", "v1.0.0", "=="))
}

func (s *HelperTestSuite) TestGetLatestPanelVersion() {
	version, err := GetLatestPanelVersion()
	s.NotEmpty(version)
	s.Nil(err)
}

func (s *HelperTestSuite) TestGetPanelVersion() {
	version, err := GetPanelVersion("v2.0.58")
	s.NotEmpty(version)
	s.Nil(err)
}
