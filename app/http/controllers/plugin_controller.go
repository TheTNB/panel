package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
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

// List
//
//	@Summary	插件列表
//	@Tags		插件
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	controllers.SuccessResponse
//	@Router		/panel/plugin/list [get]
func (r *PluginController) List(ctx http.Context) http.Response {
	plugins := r.plugin.All()
	installedPlugins, err := r.plugin.AllInstalled()
	if err != nil {
		return h.ErrorSystem(ctx)
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

	paged, total := h.Paginate(ctx, pluginArr)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// Install
//
//	@Summary	安装插件
//	@Tags		插件
//	@Produce	json
//	@Security	BearerToken
//	@Param		slug	query		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/plugin/install [post]
func (r *PluginController) Install(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")

	if err := r.plugin.Install(slug); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, "任务已提交")
}

// Uninstall
//
//	@Summary	卸载插件
//	@Tags		插件
//	@Produce	json
//	@Security	BearerToken
//	@Param		slug	query		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/plugin/uninstall [post]
func (r *PluginController) Uninstall(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")

	if err := r.plugin.Uninstall(slug); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, "任务已提交")
}

// Update
//
//	@Summary	更新插件
//	@Tags		插件
//	@Produce	json
//	@Security	BearerToken
//	@Param		slug	query		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/plugin/update [post]
func (r *PluginController) Update(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")

	if err := r.plugin.Update(slug); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, "任务已提交")
}

// UpdateShow
//
//	@Summary	更新插件首页显示状态
//	@Tags		插件
//	@Produce	json
//	@Security	BearerToken
//	@Param		slug	query		string	true	"request"
//	@Param		show	query		bool	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/plugin/updateShow [post]
func (r *PluginController) UpdateShow(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")
	show := ctx.Request().InputBool("show")

	var plugin models.Plugin
	if err := facades.Orm().Query().Where("slug", slug).First(&plugin); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "插件中心").With(map[string]any{
			"slug": slug,
			"err":  err.Error(),
		}).Info("获取插件失败")
		return h.ErrorSystem(ctx)
	}
	if plugin.ID == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "插件未安装")
	}

	plugin.Show = show
	if err := facades.Orm().Query().Save(&plugin); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "插件中心").With(map[string]any{
			"slug": slug,
			"err":  err.Error(),
		}).Info("更新插件失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, "操作成功")
}

// IsInstalled
//
//	@Summary	检查插件是否已安装
//	@Tags		插件
//	@Produce	json
//	@Security	BearerToken
//	@Param		slug	query		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/plugin/isInstalled [get]
func (r *PluginController) IsInstalled(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")

	plugin := r.plugin.GetInstalledBySlug(slug)
	info := r.plugin.GetBySlug(slug)
	if plugin.Slug != slug {
		return h.Success(ctx, http.Json{
			"name":      info.Name,
			"installed": false,
		})
	}

	return h.Success(ctx, http.Json{
		"name":      info.Name,
		"installed": true,
	})
}
