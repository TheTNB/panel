package tools

import (
	"os"
	"runtime"
)

// IsDebian 判断是否是 Debian 系统
func IsDebian() bool {
	_, err := os.Stat("/etc/debian_version")
	return err == nil
}

// IsRHEL 判断是否是 RHEL 系统
func IsRHEL() bool {
	_, err := os.Stat("/etc/redhat-release")
	return err == nil
}

// IsArm 判断是否是 ARM 架构
func IsArm() bool {
	return runtime.GOARCH == "arm" || runtime.GOARCH == "arm64"
}
