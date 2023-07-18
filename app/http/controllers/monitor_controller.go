package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"
	"panel/app/models"

	"panel/app/services"
)

type MonitorController struct {
	setting services.Setting
}

func NewMonitorController() *MonitorController {
	return &MonitorController{
		setting: services.NewSettingImpl(),
	}
}

// Switch 监控开关
func (r *MonitorController) Switch(ctx http.Context) {
	value := ctx.Request().InputBool("switch")
	err := r.setting.Set(models.SettingKeyMonitor, cast.ToString(value))
	if err != nil {
		facades.Log().Error("[面板][MonitorController] 更新监控开关失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

// SaveDays 保存监控天数
func (r *MonitorController) SaveDays(ctx http.Context) {
	days := ctx.Request().Input("days")
	err := r.setting.Set(models.SettingKeyMonitorDays, days)
	if err != nil {
		facades.Log().Error("[面板][MonitorController] 更新监控天数失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

// Clear 清空监控数据
func (r *MonitorController) Clear(ctx http.Context) {
	_, err := facades.Orm().Query().Delete(&models.Monitor{})
	if err != nil {
		facades.Log().Error("[面板][MonitorController] 清空监控数据失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

// List 监控数据列表
func (r *MonitorController) List(ctx http.Context) {
	start := ctx.Request().Input("start")
	end := ctx.Request().Input("end")
	startTime := carbon.Parse(start)
	endTime := carbon.Parse(end)

	var monitors []models.Monitor
	err := facades.Orm().Query().Where("created_at", ">=", startTime).Where("created_at", "<=", endTime).Get(&monitors)
	if err != nil {
		facades.Log().Error("[面板][MonitorController] 查询监控数据失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, monitors)
}
