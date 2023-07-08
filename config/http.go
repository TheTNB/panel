package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("http", map[string]any{
		// HTTP URL
		"url": "http://localhost",
		// HTTP Host
		"host": config.Env("APP_HOST", "0.0.0.0"),
		// HTTP Port
		"port": config.Env("APP_PORT", "8888"),
		// HTTPS Configuration
		"tls": map[string]any{
			// HTTPS Host
			"host": config.Env("APP_HOST", "0.0.0.0"),
			// HTTPS Port
			"port": config.Env("APP_PORT", "8888"),
			// SSL Certificate
			"ssl": map[string]any{
				// ca.pem
				"cert": "",
				// ca.key
				"key": "",
			},
		},
	})
}
