package str

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StrTestSuite struct {
	suite.Suite
}

func TestStrTestSuite(t *testing.T) {
	suite.Run(t, &StrTestSuite{})
}

func (s *StrTestSuite) TestPlural() {
	s.Equal("users", Plural("user"))
	s.Equal("users", Plural("users"))
}

func (s *StrTestSuite) TestSingular() {
	s.Equal("user", Singular("users"))
	s.Equal("user", Singular("user"))
}

func (s *StrTestSuite) TestSnake() {
	s.Equal("topic_comment", Snake("TopicComment"))
	s.Equal("topic_comment", Snake("topic_comment"))
}

func (s *StrTestSuite) TestCamel() {
	s.Equal("TopicComment", Camel("topic_comment"))
	s.Equal("TopicComment", Camel("TopicComment"))
}

func (s *StrTestSuite) TestLowerCamel() {
	s.Equal("topicComment", LowerCamel("topic_comment"))
	s.Equal("topicComment", LowerCamel("TopicComment"))
}
