package bootstrap

import (
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"
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

	// 初始化时区
	loc, err := time.LoadLocation(panel.Conf.MustString("app.timezone"))
	if err != nil {
		panic(fmt.Sprintf("failed to load timezone: %v", err))
	}
	time.Local = loc
	carbon.SetDefault(carbon.Default{
		Layout:       carbon.DateTimeLayout,
		Timezone:     carbon.PRC,
		WeekStartsAt: carbon.Sunday,
		Locale:       "zh-CN",
	})
}
