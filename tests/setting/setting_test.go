package setting

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/tests"
)

type SettingTestSuite struct {
	suite.Suite
	tests.TestCase
	setting internal.Setting
}

func TestSettingTestSuite(t *testing.T) {
	suite.Run(t, &SettingTestSuite{
		setting: services.NewSettingImpl(),
	})
}

func (s *SettingTestSuite) SetupTest() {

}

func (s *SettingTestSuite) TestGet() {
	a := s.setting.Get("test")
	b := s.setting.Get("test", "test")
	s.Equal("", a)
	s.Equal("test", b)
}

func (s *SettingTestSuite) TestSet() {
	err := s.setting.Set("test", "test")
	s.Nil(err)
	err = s.setting.Delete("test")
	s.Nil(err)
}

func (s *SettingTestSuite) TestDelete() {
	err := s.setting.Set("test", "test")
	s.Nil(err)
	err = s.setting.Delete("test")
	s.Nil(err)
}
