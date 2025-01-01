package systemctl

import (
	"errors"
	"fmt"

	"github.com/tnb-labs/panel/pkg/shell"
)

// Status 获取服务状态
func Status(name string) (bool, error) {
	output, err := shell.Execf("systemctl status %s | grep Active | grep -v grep | awk '{print $2}'", name)
	return output == "active", err
}

// IsEnabled 服务是否启用
func IsEnabled(name string) (bool, error) {
	out, err := shell.Execf("systemctl is-enabled '%s'", name)
	if err != nil {
		return false, fmt.Errorf("failed to check service status: %w", err)
	}

	switch out {
	case "enabled":
		return true, nil
	case "disabled":
		return false, nil
	case "masked":
		return false, errors.New("service is masked")
	case "static":
		return false, errors.New("service is statically enabled")
	case "indirect":
		return false, errors.New("service is indirectly enabled")
	default:
		return false, errors.New("unknown service status")
	}
}

// Start 启动服务
func Start(name string) error {
	_, err := shell.Execf("systemctl start %s", name)
	return err
}

// Stop 停止服务
func Stop(name string) error {
	_, err := shell.Execf("systemctl stop %s", name)
	return err
}

// Restart 重启服务
func Restart(name string) error {
	_, err := shell.Execf("systemctl restart %s", name)
	return err
}

// Reload 重载服务
func Reload(name string) error {
	_, err := shell.Execf("systemctl reload %s", name)
	return err
}

// Enable 启用服务
func Enable(name string) error {
	_, err := shell.Execf("systemctl enable %s", name)
	return err
}

// Disable 禁用服务
func Disable(name string) error {
	_, err := shell.Execf("systemctl disable %s", name)
	return err
}

// Logs 获取服务日志
func Logs(name string) (string, error) {
	return shell.Execf("journalctl -u %s", name)
}

// LogsTail 获取服务日志
func LogsTail(name string, lines int) (string, error) {
	return shell.Execf("journalctl -u %s --lines %d", name, lines)
}

// LogsClear 清空服务日志
func LogsClear(name string) error {
	if _, err := shell.Execf("journalctl --rotate -u %s", name); err != nil {
		return err
	}
	_, err := shell.Execf("journalctl --vacuum-time=1s -u %s", name)
	return err
}
