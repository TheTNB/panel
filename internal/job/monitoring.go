package job

import (
	"log/slog"
	"time"

	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/pkg/tools"
)

// Monitoring 系统监控
type Monitoring struct {
	settingRepo biz.SettingRepo
}

func NewMonitoring() *Monitoring {
	return &Monitoring{
		settingRepo: data.NewSettingRepo(),
	}
}

func (r *Monitoring) Run() {
	if app.Status != app.StatusNormal {
		return
	}

	// 将等待中的任务分发
	//task := data.NewTaskRepo()
	//_ = task.DispatchWaiting()

	monitor, err := r.settingRepo.Get(biz.SettingKeyMonitor)
	if err != nil || !cast.ToBool(monitor) {
		return
	}

	info := tools.CurrentInfo(nil, nil)

	// 去除部分数据以减少数据库存储
	info.Disk = nil
	info.Cpus = nil

	if app.Status != app.StatusNormal {
		return
	}

	if err = app.Orm.Create(&biz.Monitor{Info: info}).Error; err != nil {
		app.Logger.Warn("记录系统监控失败", slog.Any("err", err))
		return
	}

	// 删除过期数据
	dayStr, err := r.settingRepo.Get(biz.SettingKeyMonitorDays)
	if err != nil {
		return
	}
	day := cast.ToInt(dayStr)
	if day <= 0 || app.Status != app.StatusNormal {
		return
	}
	if err = app.Orm.Where("created_at < ?", time.Now().AddDate(0, 0, -day).Format(time.DateTime)).Delete(&biz.Monitor{}).Error; err != nil {
		app.Logger.Warn("删除过期系统监控失败", slog.Any("err", err))
		return
	}
}
