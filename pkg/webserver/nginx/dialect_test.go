package nginx

import (
	"strings"
	"testing"

	"github.com/libtnb/assert/check"
)

func TestDefaultSiteReuseport(t *testing.T) {
	for _, asDefault := range []bool{false, true} {
		for _, reuseport := range []bool{false, true} {
			content := defaultSiteContent(asDefault, reuseport)
			check.Equal(t, reuseport, tcpReuseportRe.MatchString(content))
			check.Equal(t, asDefault, strings.Contains(content, "default_server"))
			// quic 两条始终带 reuseport，TCP 四条随开关
			want := 2
			if reuseport {
				want = 6
			}
			check.Equal(t, want, strings.Count(content, "reuseport"))
			check.Equal(t, 2, strings.Count(content, "quic"))
		}
	}

	// 旧版安装脚本写入的默认站点监听全带 reuseport
	check.True(t, tcpReuseportRe.MatchString("    listen 80 default_server reuseport;\n    listen 443 quic default_server reuseport;\n"))
	check.False(t, tcpReuseportRe.MatchString("    listen 80 default_server;\n    listen 443 quic default_server reuseport;\n"))
}
