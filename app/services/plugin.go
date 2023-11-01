// Package services 插件服务
package services

import (
	"github.com/goravel/framework/facades"

	"panel/app/models"
	"panel/app/plugins/fail2ban"
	"panel/app/plugins/mysql57"
	"panel/app/plugins/mysql80"
	"panel/app/plugins/openresty"
	"panel/app/plugins/php74"
	"panel/app/plugins/php80"
	"panel/app/plugins/php81"
	"panel/app/plugins/php82"
	"panel/app/plugins/phpmyadmin"
	"panel/app/plugins/postgresql15"
	"panel/app/plugins/postgresql16"
	"panel/app/plugins/pureftpd"
	"panel/app/plugins/redis"
	"panel/app/plugins/s3fs"
	"panel/app/plugins/supervisor"
	"panel/app/plugins/toolbox"
)

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
		Description: openresty.Description,
		Slug:        openresty.Slug,
		Version:     openresty.Version,
		Requires:    openresty.Requires,
		Excludes:    openresty.Excludes,
		Install:     openresty.Install,
		Uninstall:   openresty.Uninstall,
		Update:      openresty.Update,
	})
	p = append(p, PanelPlugin{
		Name:        mysql57.Name,
		Description: mysql57.Description,
		Slug:        mysql57.Slug,
		Version:     mysql57.Version,
		Requires:    mysql57.Requires,
		Excludes:    mysql57.Excludes,
		Install:     mysql57.Install,
		Uninstall:   mysql57.Uninstall,
		Update:      mysql57.Update,
	})
	p = append(p, PanelPlugin{
		Name:        mysql80.Name,
		Description: mysql80.Description,
		Slug:        mysql80.Slug,
		Version:     mysql80.Version,
		Requires:    mysql80.Requires,
		Excludes:    mysql80.Excludes,
		Install:     mysql80.Install,
		Uninstall:   mysql80.Uninstall,
		Update:      mysql80.Update,
	})
	p = append(p, PanelPlugin{
		Name:        postgresql15.Name,
		Description: postgresql15.Description,
		Slug:        postgresql15.Slug,
		Version:     postgresql15.Version,
		Requires:    postgresql15.Requires,
		Excludes:    postgresql15.Excludes,
		Install:     postgresql15.Install,
		Uninstall:   postgresql15.Uninstall,
		Update:      postgresql15.Update,
	})
	p = append(p, PanelPlugin{
		Name:        postgresql16.Name,
		Description: postgresql16.Description,
		Slug:        postgresql16.Slug,
		Version:     postgresql16.Version,
		Requires:    postgresql16.Requires,
		Excludes:    postgresql16.Excludes,
		Install:     postgresql16.Install,
		Uninstall:   postgresql16.Uninstall,
		Update:      postgresql16.Update,
	})
	p = append(p, PanelPlugin{
		Name:        php74.Name,
		Description: php74.Description,
		Slug:        php74.Slug,
		Version:     php74.Version,
		Requires:    php74.Requires,
		Excludes:    php74.Excludes,
		Install:     php74.Install,
		Uninstall:   php74.Uninstall,
		Update:      php74.Update,
	})
	p = append(p, PanelPlugin{
		Name:        php80.Name,
		Description: php80.Description,
		Slug:        php80.Slug,
		Version:     php80.Version,
		Requires:    php80.Requires,
		Excludes:    php80.Excludes,
		Install:     php80.Install,
		Uninstall:   php80.Uninstall,
		Update:      php80.Update,
	})
	p = append(p, PanelPlugin{
		Name:        php81.Name,
		Description: php81.Description,
		Slug:        php81.Slug,
		Version:     php81.Version,
		Requires:    php81.Requires,
		Excludes:    php81.Excludes,
		Install:     php81.Install,
		Uninstall:   php81.Uninstall,
		Update:      php81.Update,
	})
	p = append(p, PanelPlugin{
		Name:        php82.Name,
		Description: php82.Description,
		Slug:        php82.Slug,
		Version:     php82.Version,
		Requires:    php82.Requires,
		Excludes:    php82.Excludes,
		Install:     php82.Install,
		Uninstall:   php82.Uninstall,
		Update:      php82.Update,
	})
	p = append(p, PanelPlugin{
		Name:        phpmyadmin.Name,
		Description: phpmyadmin.Description,
		Slug:        phpmyadmin.Slug,
		Version:     phpmyadmin.Version,
		Requires:    phpmyadmin.Requires,
		Excludes:    phpmyadmin.Excludes,
		Install:     phpmyadmin.Install,
		Uninstall:   phpmyadmin.Uninstall,
		Update:      phpmyadmin.Update,
	})
	p = append(p, PanelPlugin{
		Name:        pureftpd.Name,
		Description: pureftpd.Description,
		Slug:        pureftpd.Slug,
		Version:     pureftpd.Version,
		Requires:    pureftpd.Requires,
		Excludes:    pureftpd.Excludes,
		Install:     pureftpd.Install,
		Uninstall:   pureftpd.Uninstall,
		Update:      pureftpd.Update,
	})
	p = append(p, PanelPlugin{
		Name:        redis.Name,
		Description: redis.Description,
		Slug:        redis.Slug,
		Version:     redis.Version,
		Requires:    redis.Requires,
		Excludes:    redis.Excludes,
		Install:     redis.Install,
		Uninstall:   redis.Uninstall,
		Update:      redis.Update,
	})
	p = append(p, PanelPlugin{
		Name:        s3fs.Name,
		Description: s3fs.Description,
		Slug:        s3fs.Slug,
		Version:     s3fs.Version,
		Requires:    s3fs.Requires,
		Excludes:    s3fs.Excludes,
		Install:     s3fs.Install,
		Uninstall:   s3fs.Uninstall,
		Update:      s3fs.Update,
	})
	p = append(p, PanelPlugin{
		Name:        supervisor.Name,
		Description: supervisor.Description,
		Slug:        supervisor.Slug,
		Version:     supervisor.Version,
		Requires:    supervisor.Requires,
		Excludes:    supervisor.Excludes,
		Install:     supervisor.Install,
		Uninstall:   supervisor.Uninstall,
		Update:      supervisor.Update,
	})
	p = append(p, PanelPlugin{
		Name:        fail2ban.Name,
		Description: fail2ban.Description,
		Slug:        fail2ban.Slug,
		Version:     fail2ban.Version,
		Requires:    fail2ban.Requires,
		Excludes:    fail2ban.Excludes,
		Install:     fail2ban.Install,
		Uninstall:   fail2ban.Uninstall,
		Update:      fail2ban.Update,
	})
	p = append(p, PanelPlugin{
		Name:        toolbox.Name,
		Description: toolbox.Description,
		Slug:        toolbox.Slug,
		Version:     toolbox.Version,
		Requires:    toolbox.Requires,
		Excludes:    toolbox.Excludes,
		Install:     toolbox.Install,
		Uninstall:   toolbox.Uninstall,
		Update:      toolbox.Update,
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
