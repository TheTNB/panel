package services

import (
	"panel/app/models"
	"panel/packages/helper"
)

type Cron interface {
	AddToSystem(cron models.Cron)
	DeleteFromSystem(cron models.Cron)
}

type CronImpl struct {
}

func NewCronImpl() *CronImpl {
	return &CronImpl{}
}

// AddToSystem 添加到系统
func (r *CronImpl) AddToSystem(cron models.Cron) {
	if helper.IsRHEL() {
		helper.ExecShell("echo \"" + cron.Time + " " + cron.Shell + " >> " + cron.Log + " 2>&1\" >> /var/spool/cron/root")
	} else {
		helper.ExecShell("echo \"" + cron.Time + " " + cron.Shell + " >> " + cron.Log + " 2>&1\" >> /var/spool/cron/crontabs/root")
	}

	helper.ExecShell("systemctl restart crond")
}

// DeleteFromSystem 从系统中删除
func (r *CronImpl) DeleteFromSystem(cron models.Cron) {
	if helper.IsRHEL() {
		helper.ExecShell("sed -i '/" + cron.Shell + "/d' /var/spool/cron/root")
	} else {
		helper.ExecShell("sed -i '/" + cron.Shell + "/d' /var/spool/cron/crontabs/root")
	}

	helper.ExecShell("systemctl restart crond")
}
