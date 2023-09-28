package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("panel", map[string]any{
		"name":    "耗子Linux面板",
		"version": "v2.0.46",
	})
}
