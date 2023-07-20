package controllers

import (
	"sync"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/models"
	"panel/app/services"
)

type PluginController struct {
	plugin services.Plugin
	task   services.Task
}

func NewPluginController() *PluginController {
	return &PluginController{
		plugin: services.NewPluginImpl(),
		task:   services.NewTaskImpl(),
	}
}

// List 列出所有插件
func (r *PluginController) List(ctx http.Context) {
	plugins := r.plugin.All()
	installedPlugins, err := r.plugin.AllInstalled()
	if err != nil {
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	var lock sync.RWMutex
	installedPluginsMap := make(map[string]models.Plugin)

	for _, p := range installedPlugins {
		lock.Lock()
		installedPluginsMap[p.Slug] = p
		lock.Unlock()
	}

	type plugin struct {
		Name             string   `json:"name"`
		Author           string   `json:"author"`
		Description      string   `json:"description"`
		Slug             string   `json:"slug"`
		Version          string   `json:"version"`
		Requires         []string `json:"requires"`
		Excludes         []string `json:"excludes"`
		Installed        bool     `json:"installed"`
		InstalledVersion string   `json:"installed_version"`
		Show             bool     `json:"show"`
	}

	var p []plugin
	for _, item := range plugins {
		installed, installedVersion, show := false, "", false
		if _, ok := installedPluginsMap[item.Slug]; ok {
			installed = true
			installedVersion = installedPluginsMap[item.Slug].Version
			show = installedPluginsMap[item.Slug].Show
		}
		p = append(p, plugin{
			Name:             item.Name,
			Author:           item.Author,
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

	Success(ctx, p)
}

// Install 安装插件
func (r *PluginController) Install(ctx http.Context) {
	slug := ctx.Request().Input("slug")
	plugins := r.plugin.All()

	var plugin services.PanelPlugin
	check := false
	for _, item := range plugins {
		if item.Slug == slug {
			check = true
			plugin = item
			break
		}
	}
	if !check {
		Error(ctx, http.StatusBadRequest, "插件不存在")
		return
	}

	var installedPlugin models.Plugin
	if err := facades.Orm().Query().Where("slug", slug).First(&installedPlugin); err != nil {
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	if installedPlugin.ID != 0 {
		Error(ctx, http.StatusBadRequest, "插件已安装")
	}

	var task models.Task
	task.Name = "安装插件 " + plugin.Name
	task.Status = models.TaskStatusWaiting
	task.Shell = "bash /www/panel/scripts/" + plugin.Slug + "/install.sh >> /tmp/" + plugin.Slug + ".log 2>&1"
	task.Log = "/tmp/" + plugin.Slug + ".log"
	if err := facades.Orm().Query().Create(&task); err != nil {
		facades.Log().Error("[面板][PluginController] 创建任务失败: " + err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	r.task.Process(task.ID)
	Success(ctx, "任务已提交")
}

// Uninstall 卸载插件
func (r *PluginController) Uninstall(ctx http.Context) {
	slug := ctx.Request().Input("slug")
	plugins := r.plugin.All()

	var plugin services.PanelPlugin
	check := false
	for _, item := range plugins {
		if item.Slug == slug {
			check = true
			plugin = item
			break
		}
	}
	if !check {
		Error(ctx, http.StatusBadRequest, "插件不存在")
		return
	}

	var installedPlugin models.Plugin
	if err := facades.Orm().Query().Where("slug", slug).First(&installedPlugin); err != nil {
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	if installedPlugin.ID == 0 {
		Error(ctx, http.StatusBadRequest, "插件未安装")
	}

	var task models.Task
	task.Name = "卸载插件 " + plugin.Name
	task.Status = models.TaskStatusWaiting
	task.Shell = "bash /www/panel/scripts/" + plugin.Slug + "/uninstall.sh >> /tmp/" + plugin.Slug + ".log 2>&1"
	task.Log = "/tmp/" + plugin.Slug + ".log"
	if err := facades.Orm().Query().Create(&task); err != nil {
		facades.Log().Error("[面板][PluginController] 创建任务失败: " + err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	r.task.Process(task.ID)
	Success(ctx, "任务已提交")
}

// Update 更新插件
func (r *PluginController) Update(ctx http.Context) {
	slug := ctx.Request().Input("slug")
	plugins := r.plugin.All()

	var plugin services.PanelPlugin
	check := false
	for _, item := range plugins {
		if item.Slug == slug {
			check = true
			plugin = item
			break
		}
	}
	if !check {
		Error(ctx, http.StatusBadRequest, "插件不存在")
		return
	}

	var installedPlugin models.Plugin
	if err := facades.Orm().Query().Where("slug", slug).First(&installedPlugin); err != nil {
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	if installedPlugin.ID == 0 {
		Error(ctx, http.StatusBadRequest, "插件未安装")
	}

	var task models.Task
	task.Name = "更新插件 " + plugin.Name
	task.Status = models.TaskStatusWaiting
	task.Shell = "bash /www/panel/scripts/" + plugin.Slug + "/update.sh >> /tmp/" + plugin.Slug + ".log 2>&1"
	task.Log = "/tmp/" + plugin.Slug + ".log"
	if err := facades.Orm().Query().Create(&task); err != nil {
		facades.Log().Error("[面板][PluginController] 创建任务失败: " + err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	r.task.Process(task.ID)
	Success(ctx, "任务已提交")
}

// UpdateShow 更新插件首页显示状态
func (r *PluginController) UpdateShow(ctx http.Context) {
	slug := ctx.Request().Input("slug")
	show := ctx.Request().InputBool("show")

	var plugin models.Plugin
	if err := facades.Orm().Query().Where("slug", slug).First(&plugin); err != nil {
		facades.Log().Error("[面板][PluginController] 查询插件失败: " + err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	if plugin.ID == 0 {
		Error(ctx, http.StatusBadRequest, "插件未安装")
		return
	}

	plugin.Show = show
	if err := facades.Orm().Query().Save(&plugin); err != nil {
		facades.Log().Error("[面板][PluginController] 更新插件失败: " + err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, "操作成功")
}
