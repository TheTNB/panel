package user

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"panel/app/services"
	"panel/tests"
)

type UserTestSuite struct {
	suite.Suite
	tests.TestCase
	user services.User
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
}
