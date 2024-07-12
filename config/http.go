package config

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/path"
	ginfacades "github.com/goravel/gin/facades"
)

func init() {
	config := facades.Config()
	config.Add("http", map[string]any{
		// HTTP Driver
		"default": "gin",
		// HTTP Drivers
		"drivers": map[string]any{
			"gin": map[string]any{
				// Optional, default is 4096 KB
				"body_limit":   1024 * 1024 * 4,
				"header_limit": 20480,
				"route": func() (route.Route, error) {
					return ginfacades.Route("gin"), nil
				},
			},
		},
		// HTTP URL
		"url": "http://localhost",
		// HTTP Host
		"host": "",
		// HTTP Port
		"port": config.Env("APP_PORT", "8888"),
		// HTTPS Configuration
		"tls": map[string]any{
			// HTTPS Host
			"host": "",
			// HTTPS Port
			"port": config.Env("APP_PORT", "8888"),
			// SSL Certificate
			"ssl": map[string]any{
				// ca.pem
				"cert": config.Env("APP_SSL_CERT", path.Executable("storage/ssl.crt")),
				// ca.key
				"key": config.Env("APP_SSL_KEY", path.Executable("storage/ssl.key")),
			},
		},
	})
}
