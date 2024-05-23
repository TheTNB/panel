package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("panel", map[string]any{
		"name":    "耗子 Linux 面板",
		"version": "v2.2.3",
		"ssl":     config.Env("APP_SSL", false),
	})
}
