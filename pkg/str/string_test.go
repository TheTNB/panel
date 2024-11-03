package str

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StringHelperTestSuite struct {
	suite.Suite
}

func TestStringHelperTestSuite(t *testing.T) {
	suite.Run(t, &StringHelperTestSuite{})
}

func (s *StringHelperTestSuite) TestFormatBytes() {
	s.Equal("1.00 B", FormatBytes(1))
	s.Equal("1.00 KB", FormatBytes(1024))
	s.Equal("1.00 MB", FormatBytes(1024*1024))
	s.Equal("1.00 GB", FormatBytes(1024*1024*1024))
	s.Equal("1.00 TB", FormatBytes(1024*1024*1024*1024))
	s.Equal("1.00 PB", FormatBytes(1024*1024*1024*1024*1024))
	s.Equal("1.00 EB", FormatBytes(1024*1024*1024*1024*1024*1024))
	s.Equal("1.00 ZB", FormatBytes(1024*1024*1024*1024*1024*1024*1024))
	s.Equal("1.00 YB", FormatBytes(1024*1024*1024*1024*1024*1024*1024*1024))
}
