package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/models"
	"panel/app/services"
	"panel/packages/helper"
)

type MenuItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Jump  string `json:"jump"`
}

type InfoController struct {
	// Dependent services
}

func NewInfoController() *InfoController {
	return &InfoController{
		// Inject services
	}
}

func (r *InfoController) Name(ctx http.Context) {
	var setting models.Setting
	err := facades.Orm().Query().Where("key", "name").First(&setting)
	if err != nil {
		facades.Log().Error("[面板][InfoController] 查询面板名称失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, http.Json{
		"name": setting.Value,
	})
}

func (r *InfoController) Menu(ctx http.Context) {
	Success(ctx, []MenuItem{
		{Name: "home", Title: "主页", Icon: "layui-icon-home", Jump: "/"},
		{Name: "website", Title: "网站管理", Icon: "layui-icon-website", Jump: "website/list"},
		{Name: "monitor", Title: "资源监控", Icon: "layui-icon-chart-screen", Jump: "monitor"},
		{Name: "safe", Title: "系统安全", Icon: "layui-icon-auz", Jump: "safe"},
		{Name: "file", Title: "文件管理", Icon: "layui-icon-file", Jump: "file"},
		{Name: "cron", Title: "计划任务", Icon: "layui-icon-date", Jump: "cron"},
		{Name: "plugin", Title: "插件中心", Icon: "layui-icon-app", Jump: "plugin"},
		{Name: "setting", Title: "面板设置", Icon: "layui-icon-set", Jump: "setting"},
	})
}

func (r *InfoController) HomePlugins(ctx http.Context) {
	var plugins []models.Plugin
	err := facades.Orm().Query().Where("show", 1).Find(&plugins)
	if err != nil {
		facades.Log().Error("[面板][InfoController] 查询首页插件失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	type pluginsData struct {
		models.Plugin
		Name string `json:"name"`
	}

	var pluginsJson []pluginsData
	for _, plugin := range plugins {
		pluginsJson = append(pluginsJson, pluginsData{
			Plugin: plugin,
			Name:   services.NewPluginImpl().GetBySlug(plugin.Slug).Name,
		})
	}

	Success(ctx, pluginsJson)
}

func (r *InfoController) NowMonitor(ctx http.Context) {
	Success(ctx, helper.GetMonitoringInfo())
}

func (r *InfoController) SystemInfo(ctx http.Context) {
	monitorInfo := helper.GetMonitoringInfo()

	Success(ctx, http.Json{
		"os_name":       monitorInfo.Host.Platform + " " + monitorInfo.Host.PlatformVersion,
		"uptime":        fmt.Sprintf("%.2f", float64(monitorInfo.Host.Uptime)/86400),
		"panel_version": facades.Config().GetString("panel.version"),
	})
}

func (r *InfoController) InstalledDbAndPhp(ctx http.Context) {
	var php []models.Plugin
	err := facades.Orm().Query().Where("slug like ?", "php%").Find(&php)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	var mysql models.Plugin
	mysqlInstalled := true
	err = facades.Orm().Query().Where("slug like ?", "mysql%").FirstOrFail(&mysql)
	if err != nil {
		mysqlInstalled = false
	}

	var postgresql models.Plugin
	postgresqlInstalled := true
	err = facades.Orm().Query().Where("slug like ?", "postgresql%").FirstOrFail(&postgresql)
	if err != nil {
		postgresqlInstalled = false
	}

	Success(ctx, http.Json{
		"php":        php,
		"mysql":      mysqlInstalled,
		"postgresql": postgresqlInstalled,
	})
}
