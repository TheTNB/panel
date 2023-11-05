package bootstrap

import (
	"github.com/gookit/validate/locales/zhcn"
	"github.com/goravel/framework/foundation"

	"panel/config"
)

func Boot() {
	zhcn.RegisterGlobal()

	app := foundation.NewApplication()

	// Bootstrap the application
	app.Boot()

	// Bootstrap the config.
	config.Boot()
}
