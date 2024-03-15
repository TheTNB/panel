// Package services 插件服务
package services

import (
	"github.com/goravel/framework/facades"

	"panel/app/models"
	"panel/internal"
)

type PluginImpl struct {
}

func NewPluginImpl() *PluginImpl {
	return &PluginImpl{}
}

// AllInstalled 获取已安装的所有插件
func (r *PluginImpl) AllInstalled() ([]models.Plugin, error) {
	var plugins []models.Plugin
	if err := facades.Orm().Query().Get(&plugins); err != nil {
		return plugins, err
	}

	return plugins, nil
}

// All 获取所有插件
func (r *PluginImpl) All() []internal.PanelPlugin {
	var plugins = []internal.PanelPlugin{
		internal.PluginOpenResty,
		internal.PluginMySQL57,
		internal.PluginMySQL80,
		internal.PluginPostgreSQL15,
		internal.PluginPostgreSQL16,
		internal.PluginPHP74,
		internal.PluginPHP80,
		internal.PluginPHP81,
		internal.PluginPHP82,
		internal.PluginPHP83,
		internal.PluginPHPMyAdmin,
		internal.PluginPureFTPd,
		internal.PluginRedis,
		internal.PluginS3fs,
		internal.PluginRsync,
		internal.PluginSupervisor,
		internal.PluginFail2ban,
		internal.PluginToolBox,
	}

	return plugins
}

// GetBySlug 根据slug获取插件
func (r *PluginImpl) GetBySlug(slug string) internal.PanelPlugin {
	for _, item := range r.All() {
		if item.Slug == slug {
			return item
		}
	}

	return internal.PanelPlugin{}
}

// GetInstalledBySlug 根据slug获取已安装的插件
func (r *PluginImpl) GetInstalledBySlug(slug string) models.Plugin {
	var plugin models.Plugin
	if err := facades.Orm().Query().Where("slug", slug).Get(&plugin); err != nil {
		return plugin
	}

	return plugin
}
