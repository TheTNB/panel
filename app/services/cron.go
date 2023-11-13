package services

import (
	"strings"

	"panel/app/models"
	"panel/pkg/tools"
)

type Cron interface {
	AddToSystem(cron models.Cron) error
	DeleteFromSystem(cron models.Cron) error
}

type CronImpl struct {
}

func NewCronImpl() *CronImpl {
	return &CronImpl{}
}

// AddToSystem 添加到系统
func (r *CronImpl) AddToSystem(cron models.Cron) error {
	if tools.IsRHEL() {
		if _, err := tools.Exec("echo \"" + cron.Time + " " + cron.Shell + " >> " + cron.Log + " 2>&1\" >> /var/spool/cron/root"); err != nil {
			return err
		}
		if _, err := tools.Exec("systemctl restart crond"); err != nil {
			return err
		}
	} else {
		if _, err := tools.Exec("echo \"" + cron.Time + " " + cron.Shell + " >> " + cron.Log + " 2>&1\" >> /var/spool/cron/crontabs/root"); err != nil {
			return err
		}
		if _, err := tools.Exec("systemctl restart cron"); err != nil {
			return err
		}
	}

	return nil
}

// DeleteFromSystem 从系统中删除
func (r *CronImpl) DeleteFromSystem(cron models.Cron) error {
	// 需要转义Shell路径的/为\/
	cron.Shell = strings.ReplaceAll(cron.Shell, "/", "\\/")
	if tools.IsRHEL() {
		if _, err := tools.Exec("sed -i '/" + cron.Shell + "/d' /var/spool/cron/root"); err != nil {
			return err
		}
		if _, err := tools.Exec("systemctl restart crond"); err != nil {
			return err
		}
	} else {
		if _, err := tools.Exec("sed -i '/" + cron.Shell + "/d' /var/spool/cron/crontabs/root"); err != nil {
			return err
		}
		if _, err := tools.Exec("systemctl restart cron"); err != nil {
			return err
		}
	}

	return nil
}
