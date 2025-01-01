package bootstrap

import (
	"log"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/tnb-labs/panel/internal/app"
	"github.com/tnb-labs/panel/pkg/io"
)

func NewConf() (*koanf.Koanf, error) {
	config := "/usr/local/etc/panel/config.yml"
	if !io.Exists(config) {
		config = "config.yml"
	}

	conf := koanf.New(".")
	if err := conf.Load(file.Provider(config), yaml.Parser()); err != nil {
		return nil, err
	}

	initGlobal(conf)
	return conf, nil
}

func initGlobal(conf *koanf.Koanf) {
	app.Key = conf.MustString("app.key")
	if len(app.Key) != 32 {
		log.Fatalf("app key must be 32 characters")
	}

	app.Root = conf.MustString("app.root")
	app.Locale = conf.MustString("app.locale")

	// 初始化时区
	loc, err := time.LoadLocation(conf.MustString("app.timezone"))
	if err != nil {
		log.Fatalf("failed to load timezone: %v", err)
	}
	time.Local = loc
}
