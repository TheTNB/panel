package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
)

type MonitorController struct {
	setting internal.Setting
}

func NewMonitorController() *MonitorController {
	return &MonitorController{
		setting: services.NewSettingImpl(),
	}
}

// Switch 监控开关
func (r *MonitorController) Switch(ctx http.Context) http.Response {
	value := ctx.Request().InputBool("monitor")
	err := r.setting.Set(models.SettingKeyMonitor, cast.ToString(value))
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "资源监控").With(map[string]any{
			"monitor": value,
			"error":   err.Error(),
		}).Info("更新监控开关失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// SaveDays 保存监控天数
func (r *MonitorController) SaveDays(ctx http.Context) http.Response {
	days := ctx.Request().Input("days")
	err := r.setting.Set(models.SettingKeyMonitorDays, days)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "资源监控").With(map[string]any{
			"days":  days,
			"error": err.Error(),
		}).Info("更新监控开关失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// SwitchAndDays 监控开关和监控天数
func (r *MonitorController) SwitchAndDays(ctx http.Context) http.Response {
	monitor := r.setting.Get(models.SettingKeyMonitor)
	monitorDays := r.setting.Get(models.SettingKeyMonitorDays)

	return h.Success(ctx, http.Json{
		"switch": cast.ToBool(monitor),
		"days":   cast.ToInt(monitorDays),
	})
}

// Clear 清空监控数据
func (r *MonitorController) Clear(ctx http.Context) http.Response {
	_, err := facades.Orm().Query().Where("1 = 1").Delete(&models.Monitor{})
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "资源监控").With(map[string]any{
			"error": err.Error(),
		}).Info("清空监控数据失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// List 监控数据列表
func (r *MonitorController) List(ctx http.Context) http.Response {
	start := ctx.Request().InputInt64("start")
	end := ctx.Request().InputInt64("end")
	startTime := carbon.FromTimestampMilli(start)
	endTime := carbon.FromTimestampMilli(end)

	var monitors []models.Monitor
	err := facades.Orm().Query().Where("created_at >= ?", startTime.ToDateTimeString()).Where("created_at <= ?", endTime.ToDateTimeString()).Get(&monitors)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "资源监控").With(map[string]any{
			"start": startTime.ToDateTimeString(),
			"end":   endTime.ToDateTimeString(),
			"error": err.Error(),
		}).Info("获取监控数据失败")
		return h.ErrorSystem(ctx)
	}

	if len(monitors) == 0 {
		return h.Error(ctx, http.StatusNotFound, "监控数据为空")
	}

	type load struct {
		Load1  []float64 `json:"load1"`
		Load5  []float64 `json:"load5"`
		Load15 []float64 `json:"load15"`
	}
	type cpu struct {
		Percent []string `json:"percent"`
	}
	type mem struct {
		Total     string   `json:"total"`
		Available []string `json:"available"`
		Used      []string `json:"used"`
	}
	type swap struct {
		Total string   `json:"total"`
		Used  []string `json:"used"`
		Free  []string `json:"free"`
	}
	type network struct {
		Sent []string `json:"sent"`
		Recv []string `json:"recv"`
		Tx   []string `json:"tx"`
		Rx   []string `json:"rx"`
	}
	type monitorData struct {
		Times []string `json:"times"`
		Load  load     `json:"load"`
		Cpu   cpu      `json:"cpu"`
		Mem   mem      `json:"mem"`
		Swap  swap     `json:"swap"`
		Net   network  `json:"net"`
	}

	var data monitorData
	var bytesSent uint64
	var bytesRecv uint64
	var bytesSent2 uint64
	var bytesRecv2 uint64
	for _, net := range monitors[0].Info.Net {
		if net.Name == "lo" {
			continue
		}
		bytesSent += net.BytesSent
		bytesRecv += net.BytesRecv
	}
	for i, monitor := range monitors {
		// 跳过第一条数据，因为第一条数据的流量为 0
		if i == 0 {
			// MB
			data.Mem.Total = fmt.Sprintf("%.2f", float64(monitor.Info.Mem.Total)/1024/1024)
			data.Swap.Total = fmt.Sprintf("%.2f", float64(monitor.Info.Swap.Total)/1024/1024)
			continue
		}
		for _, net := range monitor.Info.Net {
			if net.Name == "lo" {
				continue
			}
			bytesSent2 += net.BytesSent
			bytesRecv2 += net.BytesRecv
		}
		data.Times = append(data.Times, monitor.CreatedAt.ToDateTimeString())
		data.Load.Load1 = append(data.Load.Load1, monitor.Info.Load.Load1)
		data.Load.Load5 = append(data.Load.Load5, monitor.Info.Load.Load5)
		data.Load.Load15 = append(data.Load.Load15, monitor.Info.Load.Load15)
		data.Cpu.Percent = append(data.Cpu.Percent, fmt.Sprintf("%.2f", monitor.Info.Percent[0]))
		data.Mem.Available = append(data.Mem.Available, fmt.Sprintf("%.2f", float64(monitor.Info.Mem.Available)/1024/1024))
		data.Mem.Used = append(data.Mem.Used, fmt.Sprintf("%.2f", float64(monitor.Info.Mem.Used)/1024/1024))
		data.Swap.Used = append(data.Swap.Used, fmt.Sprintf("%.2f", float64(monitor.Info.Swap.Used)/1024/1024))
		data.Swap.Free = append(data.Swap.Free, fmt.Sprintf("%.2f", float64(monitor.Info.Swap.Free)/1024/1024))
		data.Net.Sent = append(data.Net.Sent, fmt.Sprintf("%.2f", float64(bytesSent2/1024/1024)))
		data.Net.Recv = append(data.Net.Recv, fmt.Sprintf("%.2f", float64(bytesRecv2/1024/1024)))

		// 监控频率为 1 分钟，所以这里除以 60 即可得到每秒的流量
		data.Net.Tx = append(data.Net.Tx, fmt.Sprintf("%.2f", float64(bytesSent2-bytesSent)/60/1024/1024))
		data.Net.Rx = append(data.Net.Rx, fmt.Sprintf("%.2f", float64(bytesRecv2-bytesRecv)/60/1024/1024))

		bytesSent = bytesSent2
		bytesRecv = bytesRecv2
		bytesSent2 = 0
		bytesRecv2 = 0
	}

	return h.Success(ctx, data)
}
