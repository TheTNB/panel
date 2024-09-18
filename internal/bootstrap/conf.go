package bootstrap

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/TheTNB/panel/internal/panel"
)

func initConf() {
	panel.Conf = koanf.New(".")
	if err := panel.Conf.Load(file.Provider("config/config.yml"), yaml.Parser()); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
}

func initGlobal() {
	panel.Root = panel.Conf.MustString("app.root")
	panel.Version = panel.Conf.MustString("app.version")
	panel.Locale = panel.Conf.MustString("app.locale")
}
