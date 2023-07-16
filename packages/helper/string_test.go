package helper

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

func (s *StringHelperTestSuite) TestEmpty() {
	s.True(Empty(""))
	s.True(Empty(nil))
	s.True(Empty([]string{}))
	s.True(Empty(map[string]string{}))
	s.True(Empty(0))
	s.True(Empty(0.0))
	s.True(Empty(false))

	s.False(Empty(" "))
	s.False(Empty([]string{"Panel"}))
	s.False(Empty(map[string]string{"Panel": "HaoZi"}))
	s.False(Empty(1))
	s.False(Empty(1.0))
	s.False(Empty(true))
}

func (s *StringHelperTestSuite) TestFirstElement() {
	s.Equal("HaoZi", FirstElement([]string{"HaoZi"}))
}

func (s *StringHelperTestSuite) TestRandomNumber() {
	s.Len(RandomNumber(10), 10)
}

func (s *StringHelperTestSuite) TestRandomString() {
	s.Len(RandomString(10), 10)
}

func (s *StringHelperTestSuite) TestMD5() {
	s.Equal("e10adc3949ba59abbe56e057f20f883e", MD5("123456"))
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

func (s *StringHelperTestSuite) TestCut() {
	s.Equal("aoZ", Cut("H", "i", "HaoZi"))
}
