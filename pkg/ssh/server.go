package ssh

import (
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/acepanel/panel/v3/pkg/io"
	"github.com/acepanel/panel/v3/pkg/shell"
)

// ApplyServerPort 使已写入的 SSH 端口配置生效，保留服务和 socket 的启用状态。
func ApplyServerPort(ctx context.Context, service string, port uint) error {
	// 配置已经写入，后续应用不能因 HTTP 连接断开而半途停止。
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Minute)
	defer cancel()
	if _, err := shell.Exec(ctx, "sshd -t"); err != nil {
		return err
	}

	socket := service + ".socket"
	state, err := readServerSocket(ctx, socket)
	if err != nil {
		return err
	}
	if state["ActiveState"] != "active" && state["UnitFileState"] != "enabled" && state["UnitFileState"] != "enabled-runtime" {
		_, err = shell.Execf(ctx, "systemctl restart '%s'", service)
		return err
	}

	// Ubuntu 的 generator 在 daemon-reload 时将 sshd_config 转成 ListenStream。
	if _, err = shell.Exec(ctx, "systemctl daemon-reload"); err != nil {
		return err
	}
	state, err = readServerSocket(ctx, socket)
	if err != nil {
		return err
	}
	config, changed, err := serverSocketPortConfig(state["Listen"], port)
	if err != nil {
		return err
	}
	if changed {
		// 没有 generator 或存在静态覆盖时，只替换端口，保留 IPv4/IPv6 监听地址。
		if err = io.Write("/etc/systemd/system/"+socket+".d/zz-acepanel-port.conf", config, 0644); err != nil {
			return err
		}
		if _, err = shell.Exec(ctx, "systemctl daemon-reload"); err != nil {
			return err
		}
	}

	if state["Accept"] != "yes" {
		// Accept=no 的 sshd 持有旧监听 FD，先停止两者，避免 socket 在服务仍活跃时拒绝启动。
		if _, err = shell.Execf(ctx, "systemctl stop '%s' '%s'", socket, service); err != nil {
			return err
		}
	}
	// 由 socket 按需拉起服务；Accept=yes 的常驻服务可能与 socket 互斥。
	_, err = shell.Execf(ctx, "systemctl restart '%s'", socket)
	return err
}

func readServerSocket(ctx context.Context, socket string) (map[string]string, error) {
	output, err := shell.Execf(ctx, "systemctl show '%s' --property=LoadState,ActiveState,UnitFileState,Accept,Listen", socket)
	if err != nil {
		return nil, err
	}
	state := make(map[string]string)
	for line := range strings.SplitSeq(output, "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if key == "Listen" {
			state[key] += value + "\n"
		} else {
			state[key] = value
		}
	}
	switch state["LoadState"] {
	case "loaded", "not-found", "masked":
		return state, nil
	default:
		return nil, fmt.Errorf("failed to load %s: %s", socket, state["LoadState"])
	}
}

func serverSocketPortConfig(listen string, port uint) (string, bool, error) {
	if strings.TrimSpace(listen) == "" {
		return "", false, errors.New("SSH socket has no listening address")
	}
	newPort := strconv.FormatUint(uint64(port), 10)
	addresses := make([]string, 0)
	changed := false
	for line := range strings.SplitSeq(strings.TrimSpace(listen), "\n") {
		address, ok := strings.CutSuffix(line, " (Stream)")
		if !ok {
			return "", false, fmt.Errorf("unsupported SSH socket listener: %s", line)
		}
		host, oldPort, err := net.SplitHostPort(address)
		if err != nil {
			return "", false, err
		}
		changed = changed || oldPort != newPort
		address = net.JoinHostPort(host, newPort)
		if !slices.Contains(addresses, address) {
			addresses = append(addresses, address)
		}
	}
	var config strings.Builder
	config.WriteString("[Socket]\nListenStream=\n")
	for _, address := range addresses {
		config.WriteString("ListenStream=")
		config.WriteString(address)
		config.WriteByte('\n')
	}
	return config.String(), changed, nil
}
