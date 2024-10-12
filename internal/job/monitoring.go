package job

import (
	"time"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/pkg/types"
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

func (receiver *Monitoring) Run() {
	if types.Status != types.StatusNormal {
		return
	}

	// 将等待中的任务分发
	//task := data.NewTaskRepo()
	//_ = task.DispatchWaiting()

	monitor, err := receiver.settingRepo.Get(biz.SettingKeyMonitor)
	if err != nil || !cast.ToBool(monitor) {
		return
	}

	info := tools.GetMonitoringInfo()

	// 去除部分数据以减少数据库存储
	info.Disk = nil
	info.Cpus = nil

	if types.Status != types.StatusNormal {
		return
	}

	if err = app.Orm.Create(&biz.Monitor{Info: info}).Error; err != nil {
		app.Logger.Error("记录系统监控失败", zap.Error(err))
		return
	}

	// 删除过期数据
	dayStr, err := receiver.settingRepo.Get(biz.SettingKeyMonitorDays)
	if err != nil {
		return
	}
	day := cast.ToInt(dayStr)
	if day <= 0 || types.Status != types.StatusNormal {
		return
	}
	if err = app.Orm.Where("created_at < ?", time.Now().AddDate(0, 0, -day).Format("2006-01-02 15:04:05")).Delete(&biz.Monitor{}).Error; err != nil {
		app.Logger.Error("删除过期系统监控失败", zap.Error(err))
		return
	}
}
