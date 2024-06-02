package tools

import "fmt"

// ServiceStatus 获取服务状态
func ServiceStatus(name string) (bool, error) {
	output, err := Exec(fmt.Sprintf("systemctl status %s | grep Active | grep -v grep | awk '{print $2}'", name))
	return output == "active", err
}

// ServiceIsEnabled 服务是否启用
func ServiceIsEnabled(name string) (bool, error) {
	output, err := Exec(fmt.Sprintf("systemctl is-enabled %s", name))
	return output == "enabled", err
}

// ServiceStart 启动服务
func ServiceStart(name string) error {
	_, err := Exec(fmt.Sprintf("systemctl start %s", name))
	return err
}

// ServiceStop 停止服务
func ServiceStop(name string) error {
	_, err := Exec(fmt.Sprintf("systemctl stop %s", name))
	return err
}

// ServiceRestart 重启服务
func ServiceRestart(name string) error {
	_, err := Exec(fmt.Sprintf("systemctl restart %s", name))
	return err
}

// ServiceReload 重载服务
func ServiceReload(name string) error {
	_, err := Exec(fmt.Sprintf("systemctl reload %s", name))
	return err
}

// ServiceEnable 启用服务
func ServiceEnable(name string) error {
	_, err := Exec(fmt.Sprintf("systemctl enable %s", name))
	return err
}

// ServiceDisable 禁用服务
func ServiceDisable(name string) error {
	_, err := Exec(fmt.Sprintf("systemctl disable %s", name))
	return err
}
