package nginx

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tnb-labs/panel/pkg/io"
)

type NginxTestSuite struct {
	suite.Suite
}

func TestNginxTestSuite(t *testing.T) {
	suite.Run(t, &NginxTestSuite{})
}

func (s *NginxTestSuite) TestListen() {
	parser, err := NewParser()
	s.NoError(err)
	listen, err := parser.GetListen()
	s.NoError(err)
	s.Equal([][]string{{"80"}}, listen)
	s.NoError(parser.SetListen([][]string{{"80"}, {"443"}}))
	listen, err = parser.GetListen()
	s.NoError(err)
	s.Equal([][]string{{"80"}, {"443"}}, listen)
}

func (s *NginxTestSuite) TestServerName() {
	parser, err := NewParser()
	s.NoError(err)
	serverName, err := parser.GetServerName()
	s.NoError(err)
	s.Equal([]string{"localhost"}, serverName)
	s.NoError(parser.SetServerName([]string{"example.com"}))
	serverName, err = parser.GetServerName()
	s.NoError(err)
	s.Equal([]string{"example.com"}, serverName)
}

func (s *NginxTestSuite) TestIndex() {
	parser, err := NewParser()
	s.NoError(err)
	index, err := parser.GetIndex()
	s.NoError(err)
	s.Equal([]string{"index.php", "index.html", "index.htm"}, index)
	s.NoError(parser.SetIndex([]string{"index.html", "index.htm"}))
	index, err = parser.GetIndex()
	s.NoError(err)
	s.Equal([]string{"index.html", "index.htm"}, index)
}

func (s *NginxTestSuite) TestIndexWithComment() {
	parser, err := NewParser()
	s.NoError(err)
	index, comment, err := parser.GetIndexWithComment()
	s.NoError(err)
	s.Equal([]string{"index.php", "index.html", "index.htm"}, index)
	s.Equal([]string(nil), comment)
	s.NoError(parser.SetIndexWithComment([]string{"index.html", "index.htm"}, []string{"# 测试"}))
	index, comment, err = parser.GetIndexWithComment()
	s.NoError(err)
	s.Equal([]string{"index.html", "index.htm"}, index)
	s.Equal([]string{"# 测试"}, comment)
}

func (s *NginxTestSuite) TestRoot() {
	parser, err := NewParser()
	s.NoError(err)
	root, err := parser.GetRoot()
	s.NoError(err)
	s.Equal("/www/wwwroot/default", root)
	s.NoError(parser.SetRoot("/www/wwwroot/test"))
	root, err = parser.GetRoot()
	s.NoError(err)
	s.Equal("/www/wwwroot/test", root)
}

func (s *NginxTestSuite) TestRootWithComment() {
	parser, err := NewParser()
	s.NoError(err)
	root, comment, err := parser.GetRootWithComment()
	s.NoError(err)
	s.Equal("/www/wwwroot/default", root)
	s.Equal([]string(nil), comment)
	s.NoError(parser.SetRootWithComment("/www/wwwroot/test", []string{"# 测试"}))
	root, comment, err = parser.GetRootWithComment()
	s.NoError(err)
	s.Equal("/www/wwwroot/test", root)
	s.Equal([]string{"# 测试"}, comment)
}

func (s *NginxTestSuite) TestIncludes() {
	parser, err := NewParser()
	s.NoError(err)
	includes, comments, err := parser.GetIncludes()
	s.NoError(err)
	s.Equal([]string{"enable-php-0.conf"}, includes)
	s.Equal([][]string{[]string(nil)}, comments)
	s.NoError(parser.SetIncludes([]string{"/www/server/vhost/rewrite/default.conf"}, nil))
	includes, comments, err = parser.GetIncludes()
	s.NoError(err)
	s.Equal([]string{"/www/server/vhost/rewrite/default.conf"}, includes)
	s.Equal([][]string{[]string(nil)}, comments)
	s.NoError(parser.SetIncludes([]string{"/www/server/vhost/rewrite/test.conf"}, [][]string{{"# 伪静态规则测试"}}))
	includes, comments, err = parser.GetIncludes()
	s.NoError(err)
	s.Equal([]string{"/www/server/vhost/rewrite/test.conf"}, includes)
	s.Equal([][]string{{"# 伪静态规则测试"}}, comments)
}

func (s *NginxTestSuite) TestPHP() {
	parser, err := NewParser()
	s.NoError(err)
	s.Equal(0, parser.GetPHP())
	s.NoError(parser.SetPHP(80))
	s.Equal(80, parser.GetPHP())
	s.NoError(parser.SetPHP(0))
	s.Equal(0, parser.GetPHP())
}

func (s *NginxTestSuite) TestHTTP() {
	parser, err := NewParser()
	s.NoError(err)
	expect, err := io.Read("testdata/http.conf")
	s.NoError(err)
	s.Equal(expect, parser.Dump())
}

func (s *NginxTestSuite) TestHTTPS() {
	parser, err := NewParser()
	s.NoError(err)
	s.False(parser.GetHTTPS())
	s.NoError(parser.SetHTTPS("/www/server/vhost/cert/default.pem", "/www/server/vhost/cert/default.key"))
	s.True(parser.GetHTTPS())
	expect, err := io.Read("testdata/https.conf")
	s.NoError(err)
	s.Equal(expect, parser.Dump())
}

func (s *NginxTestSuite) TestHTTPSProtocols() {
	parser, err := NewParser()
	s.NoError(err)
	s.NoError(parser.SetHTTPS("/www/server/vhost/cert/default.pem", "/www/server/vhost/cert/default.key"))
	s.Equal([]string{"TLSv1.2", "TLSv1.3"}, parser.GetHTTPSProtocols())
	s.NoError(parser.SetHTTPSProtocols([]string{"TLSv1.3"}))
	s.Equal([]string{"TLSv1.3"}, parser.GetHTTPSProtocols())
}

func (s *NginxTestSuite) TestHTTPSCiphers() {
	parser, err := NewParser()
	s.NoError(err)
	s.NoError(parser.SetHTTPS("/www/server/vhost/cert/default.pem", "/www/server/vhost/cert/default.key"))
	s.Equal("ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-CHACHA20-POLY1305", parser.GetHTTPSCiphers())
	s.NoError(parser.SetHTTPSCiphers("TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384"))
	s.Equal("TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384", parser.GetHTTPSCiphers())
}

func (s *NginxTestSuite) TestOCSP() {
	parser, err := NewParser()
	s.NoError(err)
	s.NoError(err)
	s.NoError(parser.SetHTTPS("/www/server/vhost/cert/default.pem", "/www/server/vhost/cert/default.key"))
	s.False(parser.GetOCSP())
	s.NoError(parser.SetOCSP(false))
	s.False(parser.GetOCSP())
	s.NoError(parser.SetOCSP(true))
	s.True(parser.GetOCSP())
	s.NoError(parser.SetOCSP(false))
	s.False(parser.GetOCSP())
}

func (s *NginxTestSuite) TestHSTS() {
	parser, err := NewParser()
	s.NoError(err)
	s.NoError(parser.SetHTTPS("/www/server/vhost/cert/default.pem", "/www/server/vhost/cert/default.key"))
	s.False(parser.GetHSTS())
	s.NoError(parser.SetHSTS(false))
	s.False(parser.GetHSTS())
	s.NoError(parser.SetHSTS(true))
	s.True(parser.GetHSTS())
	s.NoError(parser.SetHSTS(false))
	s.False(parser.GetHSTS())
}

func (s *NginxTestSuite) TestHTTPSRedirect() {
	parser, err := NewParser()
	s.NoError(err)
	s.NoError(parser.SetHTTPS("/www/server/vhost/cert/default.pem", "/www/server/vhost/cert/default.key"))
	s.False(parser.GetHTTPSRedirect())
	s.NoError(parser.SetHTTPRedirect(false))
	s.False(parser.GetHTTPSRedirect())
	s.NoError(parser.SetHTTPRedirect(true))
	s.True(parser.GetHTTPSRedirect())
	s.NoError(parser.SetHTTPRedirect(false))
	s.False(parser.GetHTTPSRedirect())
}

func (s *NginxTestSuite) TestAltSvc() {
	parser, err := NewParser()
	s.NoError(err)
	s.NoError(parser.SetHTTPS("/www/server/vhost/cert/default.pem", "/www/server/vhost/cert/default.key"))
	s.Equal("", parser.GetAltSvc())
	s.NoError(parser.SetAltSvc(`'h3=":$server_port"; ma=2592000'`))
	s.Equal(`'h3=":$server_port"; ma=2592000'`, parser.GetAltSvc())
	s.NoError(parser.SetAltSvc(""))
	s.Equal("", parser.GetAltSvc())
}

func (s *NginxTestSuite) TestAccessLog() {
	parser, err := NewParser()
	s.NoError(err)
	log, err := parser.GetAccessLog()
	s.NoError(err)
	s.Equal("/www/wwwlogs/default.log", log)
	s.NoError(parser.SetAccessLog("/www/wwwlogs/access.log"))
	log, err = parser.GetAccessLog()
	s.NoError(err)
	s.Equal("/www/wwwlogs/access.log", log)
}

func (s *NginxTestSuite) TestErrorLog() {
	parser, err := NewParser()
	s.NoError(err)
	log, err := parser.GetErrorLog()
	s.NoError(err)
	s.Equal("/www/wwwlogs/default.log", log)
	s.NoError(parser.SetErrorLog("/www/wwwlogs/error.log"))
	log, err = parser.GetErrorLog()
	s.NoError(err)
	s.Equal("/www/wwwlogs/error.log", log)
}
