package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
)

type PluginService struct {
	pluginRepo biz.PluginRepo
}

func NewPluginService() *PluginService {
	return &PluginService{
		pluginRepo: data.NewPluginRepo(),
	}
}

func (s *PluginService) List(w http.ResponseWriter, r *http.Request) {
	plugins := s.pluginRepo.All()
	installedPlugins, err := s.pluginRepo.Installed()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	installedPluginsMap := make(map[string]*biz.Plugin)

	for _, p := range installedPlugins {
		installedPluginsMap[p.Slug] = p
	}

	type plugin struct {
		Name             string   `json:"name"`
		Description      string   `json:"description"`
		Slug             string   `json:"slug"`
		Version          string   `json:"version"`
		Requires         []string `json:"requires"`
		Excludes         []string `json:"excludes"`
		Installed        bool     `json:"installed"`
		InstalledVersion string   `json:"installed_version"`
		Show             bool     `json:"show"`
	}

	var pluginArr []plugin
	for _, item := range plugins {
		installed, installedVersion, show := false, "", false
		if _, ok := installedPluginsMap[item.Slug]; ok {
			installed = true
			installedVersion = installedPluginsMap[item.Slug].Version
			show = installedPluginsMap[item.Slug].Show
		}
		pluginArr = append(pluginArr, plugin{
			Name:             item.Name,
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

	paged, total := Paginate(r, pluginArr)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *PluginService) Install(w http.ResponseWriter, r *http.Request) {

}

func (s *PluginService) Uninstall(w http.ResponseWriter, r *http.Request) {

}

func (s *PluginService) Update(w http.ResponseWriter, r *http.Request) {

}

func (s *PluginService) UpdateShow(w http.ResponseWriter, r *http.Request) {

}

func (s *PluginService) IsInstalled(w http.ResponseWriter, r *http.Request) {

}
