package os

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type OSHelperTestSuite struct {
	suite.Suite
}

func TestOSHelperTestSuite(t *testing.T) {
	suite.Run(t, &OSHelperTestSuite{})
}

func (s *OSHelperTestSuite) TestIsDebian() {
	s.True(IsDebian())
}

func (s *OSHelperTestSuite) TestIsRHEL() {
	s.False(IsRHEL())
}
