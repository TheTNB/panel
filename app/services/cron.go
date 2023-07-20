package services

import (
	"strings"

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
}

// DeleteFromSystem 从系统中删除
func (r *CronImpl) DeleteFromSystem(cron models.Cron) {
	// 需要转义Shell路径的/为\/
	cron.Shell = strings.ReplaceAll(cron.Shell, "/", "\\/")
	if tools.IsRHEL() {
		tools.ExecShell("sed -i '/" + cron.Shell + "/d' /var/spool/cron/root")
	} else {
		tools.ExecShell("sed -i '/" + cron.Shell + "/d' /var/spool/cron/crontabs/root")
	}
}
