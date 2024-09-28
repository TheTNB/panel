package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/expr-lang/expr"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/pkg/api"
	"github.com/TheTNB/panel/pkg/apploader"
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
		api:       api.NewAPI(panel.Version),
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
	for app := range slices.Values(r.All()) {
		if app.Slug == slug {
			return app, nil
		}
	}
	return nil, errors.New("应用不存在")
}

func (r *appRepo) Installed() ([]*biz.App, error) {
	var apps []*biz.App
	if err := panel.Orm.Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil

}

func (r *appRepo) GetInstalled(slug string) (*biz.App, error) {
	app := new(biz.App)
	if err := panel.Orm.Where("slug = ?", slug).First(app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

func (r *appRepo) GetInstalledAll(query string, cond ...string) ([]*biz.App, error) {
	var apps []*biz.App
	if err := panel.Orm.Where(query, cond).Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *appRepo) GetHomeShow() ([]map[string]string, error) {
	var apps []*biz.App
	if err := panel.Orm.Where("show = ?", true).Order("show_order").Find(&apps).Error; err != nil {
		return nil, err
	}

	var filtered []map[string]string
	for app := range slices.Values(apps) {
		loaded, err := r.Get(app.Slug)
		if err != nil {
			continue
		}
		filtered = append(filtered, map[string]string{
			"name":        loaded.Name,
			"description": loaded.Description,
			"slug":        loaded.Slug,
			"icon":        loaded.Icon,
			"version":     app.Version,
		})
	}

	return filtered, nil
}

func (r *appRepo) IsInstalled(query string, cond ...string) (bool, error) {
	var count int64
	if len(cond) == 0 {
		if err := panel.Orm.Model(&biz.App{}).Where("slug = ?", query).Count(&count).Error; err != nil {
			return false, err
		}
	} else {
		if err := panel.Orm.Model(&biz.App{}).Where(query, cond).Count(&count).Error; err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (r *appRepo) Install(slug string) error {
	app, err := r.Get(slug)
	if err != nil {
		return err
	}

	if installed, _ := r.IsInstalled(slug); installed {
		return errors.New("应用已安装")
	}

	var shellUrl string
	for version := range slices.Values(app.Versions) {
		if version.PanelVersion == panel.Version {
			shellUrl = version.Install
			break
		}
	}
	if shellUrl == "" {
		return fmt.Errorf("应用 %s 不支持当前面板版本", app.Name)
	}

	if err = r.preCheck(app); err != nil {
		return err
	}

	task := new(biz.Task)
	task.Name = "安装应用 " + app.Name
	task.Status = biz.TaskStatusWaiting
	task.Shell = fmt.Sprintf("%s >> /tmp/%s.log 2>&1", shellUrl, app.Slug)
	task.Log = "/tmp/" + app.Slug + ".log"
	if err = r.taskRepo.Push(task); err != nil {
		return err
	}

	return err
}

func (r *appRepo) Uninstall(slug string) error {
	app, err := r.Get(slug)
	if err != nil {
		return err
	}

	if installed, _ := r.IsInstalled(slug); !installed {
		return errors.New("应用未安装")
	}

	var shellUrl string
	for version := range slices.Values(app.Versions) {
		if version.PanelVersion == panel.Version {
			shellUrl = version.Uninstall
			break
		}
	}
	if shellUrl == "" && len(app.Versions) > 0 {
		shellUrl = app.Versions[0].Uninstall
	}
	if shellUrl == "" {
		return fmt.Errorf("无法获取应用 %s 的卸载脚本", app.Name)
	}

	if err = r.preCheck(app); err != nil {
		return err
	}

	task := new(biz.Task)
	task.Name = "卸载应用 " + app.Name
	task.Status = biz.TaskStatusWaiting
	task.Shell = fmt.Sprintf("%s >> /tmp/%s.log 2>&1", shellUrl, app.Slug)
	task.Log = "/tmp/" + app.Slug + ".log"
	if err = r.taskRepo.Push(task); err != nil {
		return err
	}

	return err
}

func (r *appRepo) Update(slug string) error {
	app, err := r.Get(slug)
	if err != nil {
		return err
	}

	if installed, _ := r.IsInstalled(slug); !installed {
		return errors.New("应用未安装")
	}

	var shellUrl string
	for version := range slices.Values(app.Versions) {
		if version.PanelVersion == panel.Version {
			shellUrl = version.Update
			break
		}
	}
	if shellUrl == "" {
		return fmt.Errorf("应用 %s 不支持当前面板版本", app.Name)
	}

	if err = r.preCheck(app); err != nil {
		return err
	}

	task := new(biz.Task)
	task.Name = "更新应用 " + app.Name
	task.Status = biz.TaskStatusWaiting
	task.Shell = fmt.Sprintf("%s >> /tmp/%s.log 2>&1", shellUrl, app.Slug)
	task.Log = "/tmp/" + app.Slug + ".log"
	if err = r.taskRepo.Push(task); err != nil {
		return err
	}

	return err
}

func (r *appRepo) UpdateShow(slug string, show bool) error {
	app, err := r.GetInstalled(slug)
	if err != nil {
		return err
	}

	app.Show = show

	return panel.Orm.Save(app).Error
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
