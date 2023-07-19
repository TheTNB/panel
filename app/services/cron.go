package services

import (
	"panel/app/models"
	"panel/pkg/tools"
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
	if tools.IsRHEL() {
		tools.ExecShell("echo \"" + cron.Time + " " + cron.Shell + " >> " + cron.Log + " 2>&1\" >> /var/spool/cron/root")
	} else {
		tools.ExecShell("echo \"" + cron.Time + " " + cron.Shell + " >> " + cron.Log + " 2>&1\" >> /var/spool/cron/crontabs/root")
	}

	tools.ExecShell("systemctl restart crond")
}

// DeleteFromSystem 从系统中删除
func (r *CronImpl) DeleteFromSystem(cron models.Cron) {
	if tools.IsRHEL() {
		tools.ExecShell("sed -i '/" + cron.Shell + "/d' /var/spool/cron/root")
	} else {
		tools.ExecShell("sed -i '/" + cron.Shell + "/d' /var/spool/cron/crontabs/root")
	}

	tools.ExecShell("systemctl restart crond")
}
