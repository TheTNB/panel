package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TheTNB/panel/app/models"
	"github.com/TheTNB/panel/pkg/tools"
)

type CronImpl struct {
}

func NewCronImpl() *CronImpl {
	return &CronImpl{}
}

// AddToSystem 添加到系统
func (r *CronImpl) AddToSystem(cron models.Cron) error {
	if tools.IsRHEL() {
		if _, err := tools.Exec(fmt.Sprintf(`echo "%s %s >> %s 2>&1" >> /var/spool/cron/root`, cron.Time, cron.Shell, cron.Log)); err != nil {
			return err
		}
		return tools.ServiceRestart("crond")
	}

	if tools.IsDebian() {
		if _, err := tools.Exec(fmt.Sprintf(`echo "%s %s >> %s 2>&1" >> /var/spool/cron/crontabs/root`, cron.Time, cron.Shell, cron.Log)); err != nil {
			return err
		}
		return tools.ServiceRestart("cron")
	}

	return errors.New("不支持的系统")
}

// DeleteFromSystem 从系统中删除
func (r *CronImpl) DeleteFromSystem(cron models.Cron) error {
	// 需要转义 shell 路径的/为\/
	cron.Shell = strings.ReplaceAll(cron.Shell, "/", "\\/")
	if tools.IsRHEL() {
		if _, err := tools.Exec("sed -i '/" + cron.Shell + "/d' /var/spool/cron/root"); err != nil {
			return err
		}
		return tools.ServiceRestart("crond")
	}

	if tools.IsDebian() {
		if _, err := tools.Exec("sed -i '/" + cron.Shell + "/d' /var/spool/cron/crontabs/root"); err != nil {
			return err
		}
		return tools.ServiceRestart("cron")
	}

	return errors.New("不支持的系统")
}
