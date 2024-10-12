package bootstrap

import (
	"log"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/pkg/io"
)

func initConf() {
	config := "/usr/local/etc/panel/config.yml"
	if !io.Exists(config) {
		config = "config.yml"
	}

	app.Conf = koanf.New(".")
	if err := app.Conf.Load(file.Provider(config), yaml.Parser()); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
}

func initGlobal() {
	app.Root = app.Conf.MustString("app.root")
	app.Version = "2.3.2"
	app.Locale = app.Conf.MustString("app.locale")

	// 初始化时区
	loc, err := time.LoadLocation(app.Conf.MustString("app.timezone"))
	if err != nil {
		log.Fatalf("failed to load timezone: %v", err)
	}
	time.Local = loc
}
