package os

import (
	"os/user"

	"github.com/spf13/cast"
)

// GetUser 通过 uid 获取用户名
func GetUser(uid uint32) string {
	id := cast.ToString(uid)
	usr, err := user.LookupId(id)
	if err != nil {
		return id
	}
	return usr.Username
}

// GetGroup 通过 gid 获取组名
func GetGroup(gid uint32) string {
	id := cast.ToString(gid)
	usr, err := user.LookupGroupId(id)
	if err != nil {
		return id
	}
	return usr.Name
}
