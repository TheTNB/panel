package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TheTNB/panel/app/models"
	"github.com/TheTNB/panel/pkg/os"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type CronImpl struct {
}

func NewCronImpl() *CronImpl {
	return &CronImpl{}
}

// AddToSystem 添加到系统
func (r *CronImpl) AddToSystem(cron models.Cron) error {
	if os.IsRHEL() {
		if _, err := shell.Execf(fmt.Sprintf(`echo "%s %s >> %s 2>&1" >> /var/spool/cron/root`, cron.Time, cron.Shell, cron.Log)); err != nil {
			return err
		}
		return systemctl.Restart("crond")
	}

	if os.IsDebian() {
		if _, err := shell.Execf(fmt.Sprintf(`echo "%s %s >> %s 2>&1" >> /var/spool/cron/crontabs/root`, cron.Time, cron.Shell, cron.Log)); err != nil {
			return err
		}
		return systemctl.Restart("cron")
	}

	return errors.New("不支持的系统")
}

// DeleteFromSystem 从系统中删除
func (r *CronImpl) DeleteFromSystem(cron models.Cron) error {
	// 需要转义 shell 路径的/为\/
	cron.Shell = strings.ReplaceAll(cron.Shell, "/", "\\/")
	if os.IsRHEL() {
		if _, err := shell.Execf("sed -i '/" + cron.Shell + "/d' /var/spool/cron/root"); err != nil {
			return err
		}
		return systemctl.Restart("crond")
	}

	if os.IsDebian() {
		if _, err := shell.Execf("sed -i '/" + cron.Shell + "/d' /var/spool/cron/crontabs/root"); err != nil {
			return err
		}
		return systemctl.Restart("cron")
	}

	return errors.New("不支持的系统")
}
