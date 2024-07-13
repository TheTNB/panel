package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/console/commands"
)

type Kernel struct {
}

func (kernel *Kernel) Schedule() []schedule.Event {
	return []schedule.Event{
		facades.Schedule().Command("panel:monitoring").EveryMinute().SkipIfStillRunning(),
		facades.Schedule().Command("panel:cert-renew").DailyAt("04:00").SkipIfStillRunning(),
		facades.Schedule().Command("panel:task").DailyAt("03:30").SkipIfStillRunning(),
	}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{
		&commands.Panel{},
		&commands.Monitoring{},
		&commands.CertRenew{},
		&commands.PanelTask{},
	}
}
