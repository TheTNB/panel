package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/app/models"
	"github.com/TheTNB/panel/internal"
	"github.com/TheTNB/panel/internal/services"
)

type PluginController struct {
	plugin internal.Plugin
	task   internal.Task
}

func NewPluginController() *PluginController {
	return &PluginController{
		plugin: services.NewPluginImpl(),
		task:   services.NewTaskImpl(),
	}
}

// List 列出所有插件
func (r *PluginController) List(ctx http.Context) http.Response {
	plugins := r.plugin.All()
	installedPlugins, err := r.plugin.AllInstalled()
	if err != nil {
		return ErrorSystem(ctx)
	}

	installedPluginsMap := make(map[string]models.Plugin)

	for _, p := range installedPlugins {
		installedPluginsMap[p.Slug] = p
	}

	type plugin struct {
		Name             string   `json:"name"`
		Description      string   `json:"description"`
		Slug             string   `json:"slug"`
		Version          string   `json:"version"`
		Requires         []string `json:"requires"`
		Excludes         []string `json:"excludes"`
		Installed        bool     `json:"installed"`
		InstalledVersion string   `json:"installed_version"`
		Show             bool     `json:"show"`
	}

	var pluginArr []plugin
	for _, item := range plugins {
		installed, installedVersion, show := false, "", false
		if _, ok := installedPluginsMap[item.Slug]; ok {
			installed = true
			installedVersion = installedPluginsMap[item.Slug].Version
			show = installedPluginsMap[item.Slug].Show
		}
		pluginArr = append(pluginArr, plugin{
			Name:             item.Name,
			Description:      item.Description,
			Slug:             item.Slug,
			Version:          item.Version,
			Requires:         item.Requires,
			Excludes:         item.Excludes,
			Installed:        installed,
			InstalledVersion: installedVersion,
			Show:             show,
		})
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)
	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(pluginArr) {
		return Success(ctx, http.Json{
			"total": 0,
			"items": []plugin{},
		})
	}
	if endIndex > len(pluginArr) {
		endIndex = len(pluginArr)
	}
	pagedPlugins := pluginArr[startIndex:endIndex]

	return Success(ctx, http.Json{
		"total": len(pluginArr),
		"items": pagedPlugins,
	})
}

// Install 安装插件
func (r *PluginController) Install(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")

	if err := r.plugin.Install(slug); err != nil {
		return ErrorSystem(ctx)
	}

	return Success(ctx, "任务已提交")
}

// Uninstall 卸载插件
func (r *PluginController) Uninstall(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")

	if err := r.plugin.Uninstall(slug); err != nil {
		return ErrorSystem(ctx)
	}

	return Success(ctx, "任务已提交")
}

// Update 更新插件
func (r *PluginController) Update(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")

	if err := r.plugin.Update(slug); err != nil {
		return ErrorSystem(ctx)
	}

	return Success(ctx, "任务已提交")
}

// UpdateShow 更新插件首页显示状态
func (r *PluginController) UpdateShow(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")
	show := ctx.Request().InputBool("show")

	var plugin models.Plugin
	if err := facades.Orm().Query().Where("slug", slug).First(&plugin); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "插件中心").With(map[string]any{
			"slug": slug,
			"err":  err.Error(),
		}).Info("获取插件失败")
		return ErrorSystem(ctx)
	}
	if plugin.ID == 0 {
		return Error(ctx, http.StatusUnprocessableEntity, "插件未安装")
	}

	plugin.Show = show
	if err := facades.Orm().Query().Save(&plugin); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "插件中心").With(map[string]any{
			"slug": slug,
			"err":  err.Error(),
		}).Info("更新插件失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, "操作成功")
}
