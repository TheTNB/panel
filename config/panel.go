package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("panel", map[string]any{
		"name":     "耗子面板",
		"version":  "v2.2.27",
		"ssl":      config.Env("APP_SSL", false),
		"entrance": config.Env("APP_ENTRANCE", "/"),
	})
}
