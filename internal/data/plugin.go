package data

import (
	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/pkg/pluginloader"
	"github.com/TheTNB/panel/pkg/types"
)

type pluginRepo struct{}

func NewPluginRepo() biz.PluginRepo {
	return &pluginRepo{}
}

func (r *pluginRepo) All() []*types.Plugin {
	return pluginloader.All()
}

func (r *pluginRepo) Installed() ([]*biz.Plugin, error) {
	var plugins []*biz.Plugin
	if err := app.Orm.Find(&plugins).Error; err != nil {
		return nil, err
	}

	return plugins, nil

}

func (r *pluginRepo) Get(slug string) (*types.Plugin, error) {
	return pluginloader.Get(slug)
}

func (r *pluginRepo) GetInstalled(slug string) (*biz.Plugin, error) {
	plugin := new(biz.Plugin)
	if err := app.Orm.Where("slug = ?", slug).First(plugin).Error; err != nil {
		return nil, err
	}

	return plugin, nil
}

func (r *pluginRepo) GetInstalledAll(cond ...string) ([]*biz.Plugin, error) {
	var plugins []*biz.Plugin
	if err := app.Orm.Where(cond).Find(&plugins).Error; err != nil {
		return nil, err
	}

	return plugins, nil
}

func (r *pluginRepo) IsInstalled(cond ...string) (bool, error) {
	var count int64
	if err := app.Orm.Model(&biz.Plugin{}).Where(cond).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
