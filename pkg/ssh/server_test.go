package ssh

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/libtnb/assert/check"
	"github.com/libtnb/assert/must"
)

func TestApplyServerPort(t *testing.T) {
	tests := []struct {
		name       string
		state      string
		stop       string
		restart    string
		cancelled  bool
		queryError bool
	}{
		{name: "没有 socket 时重启常驻服务", state: "LoadState=not-found\nActiveState=inactive\n", restart: "restart ssh"},
		{name: "未启用 socket 时保留常驻服务", state: "LoadState=loaded\nActiveState=inactive\nUnitFileState=disabled\n", restart: "restart ssh"},
		{name: "socket 激活同时更新监听和常驻进程", state: "LoadState=loaded\nActiveState=active\nUnitFileState=enabled\nAccept=no\nListen=[::]:2222 (Stream)\n", stop: "stop ssh.socket ssh", restart: "restart ssh.socket"},
		{name: "临时启用的 socket 保留激活方式", state: "LoadState=loaded\nActiveState=inactive\nUnitFileState=enabled-runtime\nAccept=no\nListen=[::]:2222 (Stream)\n", stop: "stop ssh.socket ssh", restart: "restart ssh.socket"},
		{name: "按连接激活不启动互斥的常驻服务", state: "LoadState=loaded\nActiveState=active\nUnitFileState=disabled\nAccept=yes\nListen=[::]:2222 (Stream)\n", restart: "restart ssh.socket"},
		{name: "配置写入后取消请求仍完成应用", state: "LoadState=not-found\nActiveState=inactive\n", restart: "restart ssh", cancelled: true},
		{name: "查询失败时不改变服务状态", state: "LoadState=loaded\nActiveState=active\n", queryError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			logPath := filepath.Join(dir, "commands")
			must.NoError(t, os.WriteFile(filepath.Join(dir, "sshd"), []byte("#!/bin/sh\nexit 0\n"), 0755))
			must.NoError(t, os.WriteFile(filepath.Join(dir, "systemctl"), []byte(`#!/bin/sh
printf '%s\n' "$*" >> "$SSH_TEST_LOG"
if [ "$1" = show ]; then
    printf '%s' "$SSH_TEST_STATE"
    exit "$SSH_TEST_EXIT"
fi
`), 0755))
			t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
			t.Setenv("SSH_TEST_LOG", logPath)
			t.Setenv("SSH_TEST_STATE", tt.state)
			t.Setenv("SSH_TEST_EXIT", "0")
			if tt.queryError {
				t.Setenv("SSH_TEST_EXIT", "1")
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if tt.cancelled {
				cancel()
			}
			err := ApplyServerPort(ctx, "ssh", 2222)
			if tt.queryError {
				must.Error(t, err)
			} else {
				must.NoError(t, err)
			}
			log, err := os.ReadFile(logPath)
			must.NoError(t, err)
			var stops, restarts []string
			for line := range strings.SplitSeq(strings.TrimSpace(string(log)), "\n") {
				check.False(t, strings.HasPrefix(line, "enable ") || strings.HasPrefix(line, "disable "))
				if strings.HasPrefix(line, "stop ") {
					stops = append(stops, line)
				}
				if strings.HasPrefix(line, "restart ") {
					restarts = append(restarts, line)
				}
			}
			check.Equal(t, strings.Join(stops, "\n"), tt.stop)
			if tt.queryError {
				check.Len(t, restarts, 0)
			} else {
				check.Equal(t, strings.Join(restarts, "\n"), tt.restart)
			}
		})
	}
}

func TestServerSocketPortConfig(t *testing.T) {
	tests := []struct {
		name    string
		listen  string
		want    string
		changed bool
		wantErr bool
	}{
		{
			name:    "保留双栈监听地址",
			listen:  "0.0.0.0:22 (Stream)\n[::]:22 (Stream)\n",
			want:    "[Socket]\nListenStream=\nListenStream=0.0.0.0:2222\nListenStream=[::]:2222\n",
			changed: true,
		},
		{
			name:    "保留指定地址并合并重复监听",
			listen:  "127.0.0.1:22 (Stream)\n127.0.0.1:2222 (Stream)\n[fe80::1%eth0]:22 (Stream)\n",
			want:    "[Socket]\nListenStream=\nListenStream=127.0.0.1:2222\nListenStream=[fe80::1%eth0]:2222\n",
			changed: true,
		},
		{
			name:   "generator 已更新端口时无需覆盖",
			listen: "[::]:2222 (Stream)\n",
			want:   "[Socket]\nListenStream=\nListenStream=[::]:2222\n",
		},
		{name: "拒绝空监听", wantErr: true},
		{name: "拒绝 Unix socket", listen: "/run/ssh.sock (Stream)", wantErr: true},
		{name: "拒绝非 TCP 监听", listen: "0.0.0.0:22 (Datagram)", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, changed, err := serverSocketPortConfig(tt.listen, 2222)
			if tt.wantErr {
				must.Error(t, err)
				return
			}
			must.NoError(t, err)
			check.Equal(t, config, tt.want)
			check.Equal(t, changed, tt.changed)
		})
	}
}
