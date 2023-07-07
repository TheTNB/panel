package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SettingTestSuite struct {
	suite.Suite
	setting Setting
}

func TestSettingTestSuite(t *testing.T) {
	suite.Run(t, &SettingTestSuite{
		setting: NewSettingImpl(),
	})
}

func (s *SettingTestSuite) SetupTest() {

}

func (s *SettingTestSuite) TestGet() {
	/*mockOrm, mockDb, _, _ := mock.Orm()
	mockOrm.On("Query").Return(mockDb).Twice()
	mockDb.On("Where", "key", "test").Return(mockDb).Twice()
	mockDb.On("FirstOrFail", &models.Setting{}).Return(nil).Twice()
	a := s.setting.Get("test")
	b := s.setting.Get("test", "test")
	s.Equal("", a)
	s.Equal("test", b)*/
}

func (s *SettingTestSuite) TestSet() {
	/*mockOrm, mockDb, _, _ := mock.Orm()
	mockOrm.On("Query").Return(mockDb).Once()
	mockDb.On("UpdateOrCreate", &models.Setting{}, models.Setting{Key: "test"}, models.Setting{Value: "test"}).Return(nil).Once()
	err := s.setting.Set("test", "test")
	s.Nil(err)*/
}
