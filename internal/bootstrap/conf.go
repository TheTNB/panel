package bootstrap

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/TheTNB/panel/internal/app"
)

func initConf() {
	app.Conf = koanf.New(".")
	if err := app.Conf.Load(file.Provider("config/config.yml"), yaml.Parser()); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
}

func initGlobal() {
	app.Root = app.Conf.MustString("app.root")
	app.Version = app.Conf.MustString("app.version")
	app.Locale = app.Conf.MustString("app.locale")
}
