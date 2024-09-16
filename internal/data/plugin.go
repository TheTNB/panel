package data

import (
	"errors"
	"fmt"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/job"
	"github.com/TheTNB/panel/pkg/pluginloader"
	"github.com/TheTNB/panel/pkg/types"
)

type pluginRepo struct {
	taskRepo biz.TaskRepo
}

func NewPluginRepo() biz.PluginRepo {
	return &pluginRepo{
		taskRepo: NewTaskRepo(),
	}
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

func (r *pluginRepo) Install(slug string) error {
	plugin, err := r.Get(slug)
	if err != nil {
		return err
	}
	installedPlugins, err := r.Installed()
	if err != nil {
		return err
	}

	if installed, _ := r.IsInstalled("slug = ?", slug); installed {
		return errors.New("插件已安装")
	}

	pluginsMap := make(map[string]bool)

	for _, p := range installedPlugins {
		pluginsMap[p.Slug] = true
	}

	for _, require := range plugin.Requires {
		_, requireFound := pluginsMap[require]
		if !requireFound {
			return errors.New("插件 " + slug + " 需要依赖 " + require + " 插件")
		}
	}

	for _, exclude := range plugin.Excludes {
		_, excludeFound := pluginsMap[exclude]
		if excludeFound {
			return errors.New("插件 " + slug + " 不兼容 " + exclude + " 插件")
		}
	}

	task := new(biz.Task)
	task.Name = "安装插件 " + plugin.Name
	task.Status = biz.TaskStatusWaiting
	task.Shell = fmt.Sprintf("%s >> /tmp/%s.log 2>&1", plugin.Install, plugin.Slug)
	task.Log = "/tmp/" + plugin.Slug + ".log"

	if err = app.Orm.Create(task).Error; err != nil {
		return err
	}
	err = app.Queue.Push(job.NewProcessTask(r.taskRepo), []any{
		task.ID,
	})

	return err
}

func (r *pluginRepo) Uninstall(slug string) error {
	plugin, err := r.Get(slug)
	if err != nil {
		return err
	}
	installedPlugins, err := r.Installed()
	if err != nil {
		return err
	}

	if installed, _ := r.IsInstalled("slug = ?", slug); !installed {
		return errors.New("插件未安装")
	}
	pluginsMap := make(map[string]bool)

	for _, p := range installedPlugins {
		pluginsMap[p.Slug] = true
	}

	for _, require := range plugin.Requires {
		_, requireFound := pluginsMap[require]
		if !requireFound {
			return errors.New("插件 " + slug + " 需要依赖 " + require + " 插件")
		}
	}

	for _, exclude := range plugin.Excludes {
		_, excludeFound := pluginsMap[exclude]
		if excludeFound {
			return errors.New("插件 " + slug + " 不兼容 " + exclude + " 插件")
		}
	}

	task := new(biz.Task)
	task.Name = "卸载插件 " + plugin.Name
	task.Status = biz.TaskStatusWaiting
	task.Shell = fmt.Sprintf("%s >> /tmp/%s.log 2>&1", plugin.Uninstall, plugin.Slug)
	task.Log = "/tmp/" + plugin.Slug + ".log"

	if err = app.Orm.Create(task).Error; err != nil {
		return err
	}
	err = app.Queue.Push(job.NewProcessTask(r.taskRepo), []any{
		task.ID,
	})

	return err
}

func (r *pluginRepo) Update(slug string) error {
	plugin, err := r.Get(slug)
	if err != nil {
		return err
	}
	installedPlugins, err := r.Installed()
	if err != nil {
		return err
	}

	if installed, _ := r.IsInstalled("slug = ?", slug); !installed {
		return errors.New("插件未安装")
	}
	pluginsMap := make(map[string]bool)

	for _, p := range installedPlugins {
		pluginsMap[p.Slug] = true
	}

	for _, require := range plugin.Requires {
		_, requireFound := pluginsMap[require]
		if !requireFound {
			return errors.New("插件 " + slug + " 需要依赖 " + require + " 插件")
		}
	}

	for _, exclude := range plugin.Excludes {
		_, excludeFound := pluginsMap[exclude]
		if excludeFound {
			return errors.New("插件 " + slug + " 不兼容 " + exclude + " 插件")
		}
	}

	task := new(biz.Task)
	task.Name = "更新插件 " + plugin.Name
	task.Status = biz.TaskStatusWaiting
	task.Shell = fmt.Sprintf("%s >> /tmp/%s.log 2>&1", plugin.Update, plugin.Slug)
	task.Log = "/tmp/" + plugin.Slug + ".log"

	if err = app.Orm.Create(task).Error; err != nil {
		return err
	}
	err = app.Queue.Push(job.NewProcessTask(r.taskRepo), []any{
		task.ID,
	})

	return err
}
