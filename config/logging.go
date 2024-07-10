package config

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/path"
)

func init() {
	config := facades.Config()
	config.Add("logging", map[string]any{
		// Default Log Channel
		//
		// This option defines the default log channel that gets used when writing
		// messages to the logs. The name specified in this option should match
		// one of the channels defined in the "channels" configuration array.
		"default": "stack",

		// Log Channels
		//
		// Here you may configure the log channels for your application.
		// Available Drivers: "single", "daily", "custom", "stack"
		// Available Level: "debug", "info", "warning", "error", "fatal", "panic"
		"channels": map[string]any{
			"stack": map[string]any{
				"driver":   "stack",
				"channels": []string{"daily"},
			},
			"single": map[string]any{
				"driver": "single",
				"path":   path.Executable("storage/logs/panel.log"),
				"level":  "info",
				"print":  true,
			},
			"daily": map[string]any{
				"driver": "daily",
				"path":   path.Executable("storage/logs/panel.log"),
				"level":  "info",
				"days":   7,
				"print":  true,
			},
			"http": map[string]any{
				"driver": "daily",
				"path":   path.Executable("storage/logs/http.log"),
				"level":  "info",
				"days":   7,
				"print":  false,
			},
		},
	})
}
