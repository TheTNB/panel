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
		"host": "0.0.0.0",
		// HTTP Port
		"port": "8888",
		// HTTPS Configuration
		"tls": map[string]any{
			// HTTPS Host
			"host": "0.0.0.0",
			// HTTPS Port
			"port": "8899",
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
