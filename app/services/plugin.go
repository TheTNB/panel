// Package services 插件服务
package services

import (
	"github.com/goravel/framework/facades"

	"panel/app/models"
	"panel/app/plugins/mysql80"
	"panel/app/plugins/openresty"
)

// PanelPlugin 插件元数据结构
type PanelPlugin struct {
	Name        string
	Author      string
	Description string
	Slug        string
	Version     string
	Requires    []string
	Excludes    []string
}

type Plugin interface {
	AllInstalled() ([]models.Plugin, error)
	All() []PanelPlugin
}

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
func (r *PluginImpl) All() []PanelPlugin {
	var p []PanelPlugin

	p = append(p, PanelPlugin{
		Name:        openresty.Name,
		Author:      openresty.Author,
		Description: openresty.Description,
		Slug:        openresty.Slug,
		Version:     openresty.Version,
		Requires:    openresty.Requires,
		Excludes:    openresty.Excludes,
	})
	p = append(p, PanelPlugin{
		Name:        mysql80.Name,
		Author:      mysql80.Author,
		Description: mysql80.Description,
		Slug:        mysql80.Slug,
		Version:     mysql80.Version,
		Requires:    mysql80.Requires,
		Excludes:    mysql80.Excludes,
	})

	return p
}

// GetBySlug 根据slug获取插件
func (r *PluginImpl) GetBySlug(slug string) PanelPlugin {
	for _, item := range r.All() {
		if item.Slug == slug {
			return item
		}
	}

	return PanelPlugin{}
}

// GetInstalledBySlug 根据slug获取已安装的插件
func (r *PluginImpl) GetInstalledBySlug(slug string) models.Plugin {
	var plugin models.Plugin
	if err := facades.Orm().Query().Where("slug", slug).Get(&plugin); err != nil {
		return plugin
	}

	return plugin
}
