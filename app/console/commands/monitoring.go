package commands

import (
	"context"
	"strconv"

	"github.com/gookit/color"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"

	"panel/app/models"
	"panel/internal"
	"panel/internal/services"
	"panel/pkg/tools"
)

// Monitoring 系统监控
type Monitoring struct {
}

// Signature The name and signature of the console command.
func (receiver *Monitoring) Signature() string {
	return "panel:monitoring"
}

// Description The console command description.
func (receiver *Monitoring) Description() string {
	ctx := context.Background()
	return facades.Lang(ctx).Get("commands.panel:monitoring.description")
}

// Extend The console command extend.
func (receiver *Monitoring) Extend() command.Extend {
	return command.Extend{
		Category: "panel",
	}
}

// Handle Execute the console command.
func (receiver *Monitoring) Handle(console.Context) error {
	if internal.Status != internal.StatusNormal {
		return nil
	}

	// 将等待中的任务分发
	task := services.NewTaskImpl()
	_ = task.DispatchWaiting()

	setting := services.NewSettingImpl()
	monitor := setting.Get(models.SettingKeyMonitor)
	if !cast.ToBool(monitor) {
		return nil
	}

	info := tools.GetMonitoringInfo()
	translate := facades.Lang(context.Background())

	// 去除部分数据以减少数据库存储
	info.Disk = nil
	for _, cpu := range info.Cpus {
		cpu.VendorID = ""
		cpu.Family = ""
		cpu.Model = ""
		cpu.PhysicalID = ""
		cpu.ModelName = ""
		cpu.Flags = nil
	}

	if internal.Status != internal.StatusNormal {
		return nil
	}
	err := facades.Orm().Query().Create(&models.Monitor{
		Info: info,
	})
	if err != nil {
		facades.Log().Infof("[面板] 系统监控保存失败: %s", err.Error())
		color.Redf(translate.Get("commands.panel:monitoring.fail")+": %s", err.Error())
		return nil
	}

	// 删除过期数据
	monitorDays := setting.Get(models.SettingKeyMonitorDays)
	days, err := strconv.Atoi(monitorDays)
	if err != nil {
		return nil
	}

	if days <= 0 || internal.Status != internal.StatusNormal {
		return nil
	}
	if _, err = facades.Orm().Query().Where("created_at < ?", carbon.Now().SubDays(days).ToDateTimeString()).Delete(&models.Monitor{}); err != nil {
		facades.Log().Infof("[面板] 系统监控删除过期数据失败: %s", err.Error())
		return nil
	}

	return nil
}
