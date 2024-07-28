package services

import (
	"errors"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/pkg/os"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type CronImpl struct {
}

func NewCronImpl() *CronImpl {
	return &CronImpl{}
}

// AddToSystem 添加到系统
func (r *CronImpl) AddToSystem(cron models.Cron) error {
	if _, err := shell.Execf(`( crontab -l; echo "%s %s >> %s 2>&1" ) | sort - | uniq - | crontab -`, cron.Time, cron.Shell, cron.Log); err != nil {
		return err
	}

	return r.restartCron()
}

// DeleteFromSystem 从系统中删除
func (r *CronImpl) DeleteFromSystem(cron models.Cron) error {
	if _, err := shell.Execf(`( crontab -l | grep -v -F "%s %s >> %s 2>&1" ) | crontab -`, cron.Time, cron.Shell, cron.Log); err != nil {
		return err
	}

	return r.restartCron()
}

// restartCron 重启 cron 服务
func (r *CronImpl) restartCron() error {
	if os.IsRHEL() {
		return systemctl.Restart("crond")
	}

	if os.IsDebian() || os.IsUbuntu() {
		return systemctl.Restart("cron")
	}

	return errors.New("不支持的系统")
}
