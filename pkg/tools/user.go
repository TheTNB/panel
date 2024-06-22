package tools

import (
	"github.com/spf13/cast"
	"os/user"
)

// GetUser 通过 uid 获取用户名
func GetUser(uid uint32) string {
	usr, err := user.LookupId(cast.ToString(uid))
	if err != nil {
		return ""
	}
	return usr.Username
}

// GetGroup 通过 gid 获取组名
func GetGroup(gid uint32) string {
	usr, err := user.LookupGroupId(cast.ToString(gid))
	if err != nil {
		return ""
	}
	return usr.Name
}
