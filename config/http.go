package config

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
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
				"route": func() (route.Route, error) {
					return ginfacades.Route("gin"), nil
				},
			},
		},
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
