package user

import (
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"github.com/TheTNB/panel/app/models"
	"github.com/TheTNB/panel/internal"
	"github.com/TheTNB/panel/internal/services"
	"github.com/TheTNB/panel/tests"
)

type UserTestSuite struct {
	suite.Suite
	tests.TestCase
	user internal.User
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, &UserTestSuite{
		user: services.NewUserImpl(),
	})
}

func (s *UserTestSuite) SetupTest() {

}

func (s *UserTestSuite) TestCreate() {
	user, err := s.user.Create("haozi", "123456")
	s.Nil(err)
	s.Equal("haozi", user.Username)
	_, err = facades.Orm().Query().Where("username", "haozi").Delete(&models.User{})
	s.Nil(err)
}

func (s *UserTestSuite) TestUpdate() {
	user, err := s.user.Create("haozi", "123456")
	s.Nil(err)
	s.Equal("haozi", user.Username)
	user.Username = "haozi2"
	user, err = s.user.Update(user)
	s.Nil(err)
	s.Equal("haozi2", user.Username)
	_, err = facades.Orm().Query().Where("username", "haozi").Delete(&models.User{})
	s.Nil(err)
}
