package os

import (
	"os"
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

// IsUbuntu 判断是否是 Ubuntu 系统
func IsUbuntu() bool {
	_, err := os.Stat("/etc/lsb-release")
	return err == nil
}
