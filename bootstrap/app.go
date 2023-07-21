package bootstrap

import (
	"github.com/goravel/framework/foundation"
	"github.com/goravel/framework/support/carbon"

	"panel/config"
)

func Boot() {
	app := foundation.NewApplication()

	// Bootstrap the application
	app.Boot()

	// Bootstrap the config.
	config.Boot()

	// 设置 Carbon 时区
	carbon.SetTimezone(carbon.PRC)
}
