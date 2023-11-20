package tools

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
	if IsWindows() {
		return
	}
	s.True(IsDebian())
}

func (s *OSHelperTestSuite) TestIsRHEL() {
	if IsWindows() {
		return
	}
	s.False(IsRHEL())
}

func (s *OSHelperTestSuite) TestIsArm() {
	s.False(IsArm())
}
