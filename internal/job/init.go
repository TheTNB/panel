package job

import (
	"github.com/robfig/cron/v3"
)

func Boot(c *cron.Cron) error {
	if _, err := c.AddJob("* * * * *", NewMonitoring()); err != nil {
		return err
	}
	if _, err := c.AddJob("0 4 * * *", NewCertRenew()); err != nil {
		return err
	}
	if _, err := c.AddJob("0 2 * * *", NewPanelTask()); err != nil {
		return err
	}

	return nil
}
