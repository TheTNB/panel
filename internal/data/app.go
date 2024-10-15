package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/expr-lang/expr"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/pkg/api"
	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/str"
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
		api:       api.NewAPI(app.Version),
	}
}

func (r *appRepo) All() api.Apps {
	cached, err := r.cacheRepo.Get(biz.CacheKeyApps)
	if err != nil {
		return nil
	}
	var apps api.Apps
	if err = json.Unmarshal([]byte(cached), &apps); err != nil {
		return nil
	}
	return apps
}

func (r *appRepo) Get(slug string) (*api.App, error) {
	for item := range slices.Values(r.All()) {
		if item.Slug == slug {
			return item, nil
		}
	}
	return nil, errors.New("应用不存在")
}

func (r *appRepo) UpdateExist(slug string) bool {
	item, err := r.Get(slug)
	if err != nil {
		return false
	}
	installed, err := r.GetInstalled(slug)
	if err != nil {
		return false
	}

	for channel := range slices.Values(item.Channels) {
		if channel.Slug == installed.Channel {
			current := str.FirstElement(channel.Subs)
			if current != nil && current.Version != installed.Version {
				return true
			}
		}
	}

	return false
}

func (r *appRepo) Installed() ([]*biz.App, error) {
	var apps []*biz.App
	if err := app.Orm.Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil

}

func (r *appRepo) GetInstalled(slug string) (*biz.App, error) {
	installed := new(biz.App)
	if err := app.Orm.Where("slug = ?", slug).First(installed).Error; err != nil {
		return nil, err
	}

	return installed, nil
}

func (r *appRepo) GetInstalledAll(query string, cond ...string) ([]*biz.App, error) {
	var apps []*biz.App
	if err := app.Orm.Where(query, cond).Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *appRepo) GetHomeShow() ([]map[string]string, error) {
	var apps []*biz.App
	if err := app.Orm.Where("show = ?", true).Order("show_order").Find(&apps).Error; err != nil {
		return nil, err
	}

	var filtered []map[string]string
	for item := range slices.Values(apps) {
		loaded, err := r.Get(item.Slug)
		if err != nil {
			continue
		}
		filtered = append(filtered, map[string]string{
			"name":        loaded.Name,
			"description": loaded.Description,
			"slug":        loaded.Slug,
			"icon":        loaded.Icon,
			"version":     item.Version,
		})
	}

	return filtered, nil
}

func (r *appRepo) IsInstalled(query string, cond ...string) (bool, error) {
	var count int64
	if len(cond) == 0 {
		if err := app.Orm.Model(&biz.App{}).Where("slug = ?", query).Count(&count).Error; err != nil {
			return false, err
		}
	} else {
		if err := app.Orm.Model(&biz.App{}).Where(query, cond).Count(&count).Error; err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (r *appRepo) Install(channel, slug string) error {
	item, err := r.Get(slug)
	if err != nil {
		return err
	}
	panel, err := version.NewVersion(app.Version)
	if err != nil {
		return err
	}

	if installed, _ := r.IsInstalled(slug); installed {
		return errors.New("应用已安装")
	}

	shellUrl, shellChannel, shellVersion := "", "", ""
	for ch := range slices.Values(item.Channels) {
		vs, err := version.NewVersion(ch.Panel)
		if err != nil {
			continue
		}
		if ch.Slug == channel {
			if vs.GreaterThan(panel) {
				return fmt.Errorf("应用 %s 需要面板版本 %s，当前版本 %s", item.Name, ch.Panel, app.Version)
			}
			shellUrl = ch.Install
			shellChannel = ch.Slug
			shellVersion = str.FirstElement(ch.Subs).Version
			break
		}
	}
	if shellUrl == "" {
		return fmt.Errorf("应用 %s 不支持当前面板版本", item.Name)
	}

	if err = r.preCheck(item); err != nil {
		return err
	}

	task := new(biz.Task)
	task.Name = "安装应用 " + item.Name
	task.Status = biz.TaskStatusWaiting
	task.Shell = fmt.Sprintf(`curl -fsLm 10 --retry 3 "%s" | bash -s -- "%s" "%s" >> /tmp/%s.log 2>&1`, shellUrl, shellChannel, shellVersion, item.Slug)
	task.Log = "/tmp/" + item.Slug + ".log"
	if err = r.taskRepo.Push(task); err != nil {
		return err
	}

	return err
}

func (r *appRepo) UnInstall(slug string) error {
	item, err := r.Get(slug)
	if err != nil {
		return err
	}
	panel, err := version.NewVersion(app.Version)
	if err != nil {
		return err
	}

	if installed, _ := r.IsInstalled(slug); !installed {
		return errors.New("应用未安装")
	}
	installed, err := r.GetInstalled(slug)
	if err != nil {
		return err
	}

	shellUrl, shellChannel, shellVersion := "", "", ""
	for ch := range slices.Values(item.Channels) {
		vs, err := version.NewVersion(ch.Panel)
		if err != nil {
			continue
		}
		if ch.Slug == installed.Channel {
			if vs.GreaterThan(panel) {
				return fmt.Errorf("应用 %s 需要面板版本 %s，当前版本 %s", item.Name, ch.Panel, app.Version)
			}
			shellUrl = ch.Uninstall
			shellChannel = ch.Slug
			shellVersion = installed.Version
			break
		}
	}
	if shellUrl == "" {
		return fmt.Errorf("无法获取应用 %s 的卸载脚本", item.Name)
	}

	if err = r.preCheck(item); err != nil {
		return err
	}

	task := new(biz.Task)
	task.Name = "卸载应用 " + item.Name
	task.Status = biz.TaskStatusWaiting
	task.Shell = fmt.Sprintf(`curl -fsLm 10 --retry 3 "%s" | bash -s -- "%s" "%s" >> /tmp/%s.log 2>&1`, shellUrl, shellChannel, shellVersion, item.Slug)
	task.Log = "/tmp/" + item.Slug + ".log"
	if err = r.taskRepo.Push(task); err != nil {
		return err
	}

	return err
}

func (r *appRepo) Update(slug string) error {
	item, err := r.Get(slug)
	if err != nil {
		return err
	}
	panel, err := version.NewVersion(app.Version)
	if err != nil {
		return err
	}

	if installed, _ := r.IsInstalled(slug); !installed {
		return errors.New("应用未安装")
	}
	installed, err := r.GetInstalled(slug)
	if err != nil {
		return err
	}

	shellUrl, shellChannel, shellVersion := "", "", ""
	for ch := range slices.Values(item.Channels) {
		vs, err := version.NewVersion(ch.Panel)
		if err != nil {
			continue
		}
		if ch.Slug == installed.Channel {
			if vs.GreaterThan(panel) {
				return fmt.Errorf("应用 %s 需要面板版本 %s，当前版本 %s", item.Name, ch.Panel, app.Version)
			}
			shellUrl = ch.Update
			shellChannel = ch.Slug
			shellVersion = str.FirstElement(ch.Subs).Version
			break
		}
	}
	if shellUrl == "" {
		return fmt.Errorf("应用 %s 不支持当前面板版本", item.Name)
	}

	if err = r.preCheck(item); err != nil {
		return err
	}

	task := new(biz.Task)
	task.Name = "升级应用 " + item.Name
	task.Status = biz.TaskStatusWaiting
	task.Shell = fmt.Sprintf(`curl -fsLm 10 --retry 3 "%s" | bash -s -- "%s" "%s" >> /tmp/%s.log 2>&1`, shellUrl, shellChannel, shellVersion, item.Slug)
	task.Log = "/tmp/" + item.Slug + ".log"
	if err = r.taskRepo.Push(task); err != nil {
		return err
	}

	return err
}

func (r *appRepo) UpdateShow(slug string, show bool) error {
	item, err := r.GetInstalled(slug)
	if err != nil {
		return err
	}

	item.Show = show

	return app.Orm.Save(item).Error
}

func (r *appRepo) UpdateCache() error {
	remote, err := r.api.Apps()
	if err != nil {
		return err
	}

	// 去除本地不存在的应用
	*remote = slices.Clip(slices.DeleteFunc(*remote, func(app *api.App) bool {
		_, err = apploader.Get(app.Slug)
		return err != nil
	}))

	encoded, err := json.Marshal(remote)
	if err != nil {
		return err
	}

	return r.cacheRepo.Set(biz.CacheKeyApps, string(encoded))
}

func (r *appRepo) preCheck(app *api.App) error {
	var apps []string
	var installed []string

	all := r.All()
	for _, item := range all {
		apps = append(apps, item.Slug)
	}
	installedApps, err := r.Installed()
	if err != nil {
		return err
	}
	for _, item := range installedApps {
		installed = append(installed, item.Slug)
	}

	env := map[string]any{
		"apps":      apps,
		"installed": installed,
	}
	output, err := expr.Eval(app.Depends, env)
	if err != nil {
		return err
	}

	result := cast.ToString(output)
	if result != "ok" {
		return fmt.Errorf("应用 %s %s", app.Name, result)
	}

	return nil
}
