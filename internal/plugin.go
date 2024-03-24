package internal

import "panel/app/models"

// PanelPlugin 插件元数据结构
type PanelPlugin struct {
	Name        string
	Description string
	Slug        string
	Version     string
	Requires    []string
	Excludes    []string
	Install     string
	Uninstall   string
	Update      string
}

type Plugin interface {
	AllInstalled() ([]models.Plugin, error)
	All() []PanelPlugin
	GetBySlug(slug string) PanelPlugin
	GetInstalledBySlug(slug string) models.Plugin
	Install(slug string) error
	Uninstall(slug string) error
	Update(slug string) error
}
