package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("queue", map[string]any{
		// Default Queue Connection Name
		"default": "sync",

		// Queue Connections
		//
		// Here you may configure the connection information for each server that is used by your application.
		// Drivers: "sync", "redis"
		"connections": map[string]any{
			"sync": map[string]any{
				"driver": "sync",
			},
		},
	})
}
