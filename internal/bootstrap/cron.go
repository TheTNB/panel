package bootstrap

import (
	"fmt"

	"github.com/robfig/cron/v3"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/job"
	pkgcron "github.com/TheTNB/panel/pkg/cron"
)

func initCron() {
	logger := pkgcron.NewLogger(app.Logger, app.Conf.Bool("app.debug"))
	c := cron.New(
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)),
		cron.WithLogger(logger),
		cron.WithChain(cron.Recover(logger), cron.SkipIfStillRunning(logger)),
	)
	app.Cron = c

	if err := job.Boot(app.Cron); err != nil {
		panic(fmt.Sprintf("failed to boot cron jobs: %v", err))
	}

	c.Start()
}
