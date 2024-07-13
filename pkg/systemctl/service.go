package systemctl

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/TheTNB/panel/v2/pkg/shell"
)

// Status 获取服务状态
func Status(name string) (bool, error) {
	output, err := shell.Execf("systemctl status %s | grep Active | grep -v grep | awk '{print $2}'", name)
	return output == "active", err
}

// IsEnabled 服务是否启用
func IsEnabled(name string) (bool, error) {
	cmd := exec.Command("systemctl", "is-enabled", name)
	output, _ := cmd.CombinedOutput()
	status := strings.TrimSpace(string(output))

	switch status {
	case "enabled":
		return true, nil
	case "disabled":
		return false, nil
	case "masked":
		return false, errors.New("服务已被屏蔽")
	case "static":
		return false, errors.New("服务已被静态启用")
	case "indirect":
		return false, errors.New("服务已被间接启用")
	default:
		return false, errors.New("无法确定服务状态")
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
