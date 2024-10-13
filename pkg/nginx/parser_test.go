package nginx

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type NginxTestSuite struct {
	parser *Parser
	suite.Suite
}

func TestNginxTestSuite(t *testing.T) {
	parser, err := NewParser()
	if err != nil {
		t.Errorf("parse error %v", err)
	}
	suite.Run(t, &NginxTestSuite{
		parser: parser,
	})
}

func (suite *NginxTestSuite) TestA() {
	suite.NoError(suite.parser.SetPHP(81))
}
