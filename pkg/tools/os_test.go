package tools

import (
	"testing"

	"github.com/goravel/framework/support/env"
	"github.com/stretchr/testify/suite"
)

type OSHelperTestSuite struct {
	suite.Suite
}

func TestOSHelperTestSuite(t *testing.T) {
	suite.Run(t, &OSHelperTestSuite{})
}

func (s *OSHelperTestSuite) TestIsDebian() {
	if env.IsWindows() {
		return
	}
	s.True(IsDebian())
}

func (s *OSHelperTestSuite) TestIsRHEL() {
	if env.IsWindows() {
		return
	}
	s.False(IsRHEL())
}
