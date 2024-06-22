package tools

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/TheTNB/panel/pkg/shell"
)

// ServiceStatus 获取服务状态
func ServiceStatus(name string) (bool, error) {
	output, err := shell.Execf("systemctl status %s | grep Active | grep -v grep | awk '{print $2}'", name)
	return output == "active", err
}

// ServiceIsEnabled 服务是否启用
func ServiceIsEnabled(name string) (bool, error) {
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

// ServiceStart 启动服务
func ServiceStart(name string) error {
	_, err := shell.Execf("systemctl start %s", name)
	return err
}

// ServiceStop 停止服务
func ServiceStop(name string) error {
	_, err := shell.Execf("systemctl stop %s", name)
	return err
}

// ServiceRestart 重启服务
func ServiceRestart(name string) error {
	_, err := shell.Execf("systemctl restart %s", name)
	return err
}

// ServiceReload 重载服务
func ServiceReload(name string) error {
	_, err := shell.Execf("systemctl reload %s", name)
	return err
}

// ServiceEnable 启用服务
func ServiceEnable(name string) error {
	_, err := shell.Execf("systemctl enable %s", name)
	return err
}

// ServiceDisable 禁用服务
func ServiceDisable(name string) error {
	_, err := shell.Execf("systemctl disable %s", name)
	return err
}
