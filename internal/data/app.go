package data

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/job"
	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/pkg/api"
	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

type appRepo struct {
	cacheRepo biz.CacheRepo
	taskRepo  biz.TaskRepo
	api       *api.API
}

func NewAppRepo() biz.AppRepo {
	return &appRepo{
		cacheRepo: NewCacheRepo(),
		taskRepo:  NewTaskRepo(),
	}
}

func (r *appRepo) getCached() ([]*types.App, error) {
	cached, err := r.cacheRepo.Get(biz.CacheKeyApps)
	if err != nil {
		return nil, err
	}
	var apps []*types.App
	if err = json.Unmarshal([]byte(cached), &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *appRepo) All() []*types.App {
	return apploader.All()
}

func (r *appRepo) Installed() ([]*biz.App, error) {
	var plugins []*biz.App
	if err := panel.Orm.Find(&plugins).Error; err != nil {
		return nil, err
	}

	return plugins, nil

}

func (r *appRepo) Get(slug string) (*types.App, error) {
	return apploader.Get(slug)
}

func (r *appRepo) GetInstalled(slug string) (*biz.App, error) {
	plugin := new(biz.App)
	if err := panel.Orm.Where("slug = ?", slug).First(plugin).Error; err != nil {
		return nil, err
	}

	return plugin, nil
}

func (r *appRepo) GetInstalledAll(cond ...string) ([]*biz.App, error) {
	var plugins []*biz.App
	if err := panel.Orm.Where(cond).Find(&plugins).Error; err != nil {
		return nil, err
	}

	return plugins, nil
}

func (r *appRepo) IsInstalled(cond ...string) (bool, error) {
	var count int64
	if err := panel.Orm.Model(&biz.App{}).Where(cond).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *appRepo) Install(slug string) error {
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

	if err = panel.Orm.Create(task).Error; err != nil {
		return err
	}
	err = panel.Queue.Push(job.NewProcessTask(r.taskRepo), []any{
		task.ID,
	})

	return err
}

func (r *appRepo) Uninstall(slug string) error {
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

	if err = panel.Orm.Create(task).Error; err != nil {
		return err
	}
	err = panel.Queue.Push(job.NewProcessTask(r.taskRepo), []any{
		task.ID,
	})

	return err
}

func (r *appRepo) Update(slug string) error {
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

	if err = panel.Orm.Create(task).Error; err != nil {
		return err
	}
	err = panel.Queue.Push(job.NewProcessTask(r.taskRepo), []any{
		task.ID,
	})

	return err
}

func (r *appRepo) UpdateShow(slug string, show bool) error {
	plugin, err := r.GetInstalled(slug)
	if err != nil {
		return err
	}

	plugin.Show = show

	return panel.Orm.Save(plugin).Error
}
