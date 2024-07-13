package bootstrap

import (
	"runtime/debug"

	"github.com/gookit/validate/locales/zhcn"
	"github.com/goravel/framework/foundation"

	"github.com/TheTNB/panel/v2/config"
)

func Boot() {
	debug.SetGCPercent(10)
	debug.SetMemoryLimit(64 << 20)

	zhcn.RegisterGlobal()

	app := foundation.NewApplication()

	// Bootstrap the application
	app.Boot()

	// Bootstrap the config.
	config.Boot()
}
