package controllers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/models"
	"panel/app/services"
	"panel/pkg/tools"
)

type MenuItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Jump  string `json:"jump"`
}

type InfoController struct {
	plugin services.Plugin
}

func NewInfoController() *InfoController {
	return &InfoController{
		plugin: services.NewPluginImpl(),
	}
}

// Name 获取面板名称
func (c *InfoController) Name(ctx http.Context) http.Response {
	var setting models.Setting
	err := facades.Orm().Query().Where("key", "name").First(&setting)
	if err != nil {
		facades.Log().Error("[面板][InfoController] 查询面板名称失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, http.Json{
		"name": setting.Value,
	})
}

// Menu 获取面板菜单
func (c *InfoController) Menu(ctx http.Context) http.Response {
	return Success(ctx, []MenuItem{
		{Name: "home", Title: "主页", Icon: "layui-icon-home", Jump: "/"},
		{Name: "website", Title: "网站管理", Icon: "layui-icon-website", Jump: "website/list"},
		{Name: "monitor", Title: "资源监控", Icon: "layui-icon-chart-screen", Jump: "monitor"},
		{Name: "safe", Title: "系统安全", Icon: "layui-icon-auz", Jump: "safe"},
		/*{Name: "file", Title: "文件管理", Icon: "layui-icon-file", Jump: "file"},*/
		{Name: "cron", Title: "计划任务", Icon: "layui-icon-date", Jump: "cron"},
		{Name: "ssh", Title: "SSH", Icon: "layui-icon-layer", Jump: "ssh"},
		{Name: "plugin", Title: "插件中心", Icon: "layui-icon-app", Jump: "plugin"},
		{Name: "setting", Title: "面板设置", Icon: "layui-icon-set", Jump: "setting"},
	})
}

// HomePlugins 获取首页插件
func (c *InfoController) HomePlugins(ctx http.Context) http.Response {
	var plugins []models.Plugin
	err := facades.Orm().Query().Where("show", 1).Find(&plugins)
	if err != nil {
		facades.Log().Error("[面板][InfoController] 查询首页插件失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
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

	return Success(ctx, pluginsJson)
}

// NowMonitor 获取当前监控信息
func (c *InfoController) NowMonitor(ctx http.Context) http.Response {
	return Success(ctx, tools.GetMonitoringInfo())
}

// SystemInfo 获取系统信息
func (c *InfoController) SystemInfo(ctx http.Context) http.Response {
	monitorInfo := tools.GetMonitoringInfo()

	return Success(ctx, http.Json{
		"os_name":       monitorInfo.Host.Platform + " " + monitorInfo.Host.PlatformVersion,
		"uptime":        fmt.Sprintf("%.2f", float64(monitorInfo.Host.Uptime)/86400),
		"panel_version": facades.Config().GetString("panel.version"),
	})
}

// InstalledDbAndPhp 获取已安装的数据库和 PHP 版本
func (c *InfoController) InstalledDbAndPhp(ctx http.Context) http.Response {
	var php []models.Plugin
	err := facades.Orm().Query().Where("slug like ?", "php%").Find(&php)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
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

	type data struct {
		Slug string `json:"slug"`
		Name string `json:"name"`
	}
	var phpData []data
	phpData = append(phpData, data{Slug: "0", Name: "不使用"})
	for _, p := range php {
		match := regexp.MustCompile(`php(\d+)`).FindStringSubmatch(p.Slug)
		if len(match) == 0 {
			continue
		}

		phpData = append(phpData, data{Slug: strings.ReplaceAll(p.Slug, "php", ""), Name: c.plugin.GetBySlug(p.Slug).Name})
	}

	return Success(ctx, http.Json{
		"php":        phpData,
		"mysql":      mysqlInstalled,
		"postgresql": postgresqlInstalled,
	})
}

// CheckUpdate 检查面板更新
func (c *InfoController) CheckUpdate(ctx http.Context) http.Response {
	version := facades.Config().GetString("panel.version")
	remote, err := tools.GetLatestPanelVersion()
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, "获取最新版本失败")
	}

	if version == remote.Version {
		return Success(ctx, http.Json{
			"update":  false,
			"version": remote.Version,
			"name":    remote.Name,
			"body":    remote.Body,
			"date":    remote.Date,
		})
	}

	return Success(ctx, http.Json{
		"update":  true,
		"version": remote.Version,
		"name":    remote.Name,
		"body":    remote.Body,
		"date":    remote.Date,
	})
}

// Update 更新面板
func (c *InfoController) Update(ctx http.Context) http.Response {
	var task models.Task
	err := facades.Orm().Query().Where("status", models.TaskStatusRunning).OrWhere("status", models.TaskStatusWaiting).FirstOrFail(&task)
	if err == nil {
		return Error(ctx, http.StatusInternalServerError, "当前有任务正在执行，禁止更新")
	}

	err = tools.UpdatePanel()
	if err != nil {
		facades.Log().Error("[面板][InfoController] 更新面板失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
	}

	return Success(ctx, nil)
}

// Restart 重启面板
func (c *InfoController) Restart(ctx http.Context) http.Response {
	var task models.Task
	err := facades.Orm().Query().Where("status", models.TaskStatusRunning).OrWhere("status", models.TaskStatusWaiting).FirstOrFail(&task)
	if err == nil {
		return Error(ctx, http.StatusInternalServerError, "当前有任务正在执行，禁止重启")
	}

	tools.Exec("systemctl restart panel")
	return Success(ctx, nil)
}
