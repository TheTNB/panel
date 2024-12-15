package bootstrap

import (
	"log/slog"

	"github.com/knadh/koanf/v2"
	"github.com/robfig/cron/v3"

	pkgcron "github.com/TheTNB/panel/pkg/cron"
)

func NewCron(conf *koanf.Koanf, log *slog.Logger) *cron.Cron {
	logger := pkgcron.NewLogger(log, conf.Bool("app.debug"))

	return cron.New(
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)),
		cron.WithLogger(logger),
		cron.WithChain(cron.Recover(logger), cron.SkipIfStillRunning(logger)),
	)
}
