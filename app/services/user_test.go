package services

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/goravel/framework/testing/mock"

	"panel/app/models"
)

type UserTestSuite struct {
	suite.Suite
	user User
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, &UserTestSuite{
		user: NewUserImpl(),
	})
}

func (s *UserTestSuite) SetupTest() {

}

func (s *UserTestSuite) TestCreate() {
	mockOrm, mockDb, _, _ := mock.Orm()
	mockOrm.On("Query").Return(mockDb).Once()
	mockDb.On("Create", &models.User{
		Username: "haozi",
		Password: "123456",
	}).Return(nil).Once()
	user, err := s.user.Create("haozi", "123456")
	s.Nil(err)
	s.Equal("haozi", user.Username)
	mockOrm.AssertExpectations(s.T())
	mockDb.AssertExpectations(s.T())
}
