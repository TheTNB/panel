package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/types"
)

type AppService struct {
	appRepo biz.AppRepo
}

func NewAppService() *AppService {
	return &AppService{
		appRepo: data.NewAppRepo(),
	}
}

func (s *AppService) List(w http.ResponseWriter, r *http.Request) {
	all := s.appRepo.All()
	installedApps, err := s.appRepo.Installed()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	installedAppMap := make(map[string]*biz.App)

	for _, p := range installedApps {
		installedAppMap[p.Slug] = p
	}

	var apps []types.AppCenter
	for _, item := range all {
		installed, installedChannel, installedVersion, updateExist, show := false, "", "", false, false
		if _, ok := installedAppMap[item.Slug]; ok {
			installed = true
			installedChannel = installedAppMap[item.Slug].Channel
			installedVersion = installedAppMap[item.Slug].Version
			updateExist = s.appRepo.UpdateExist(item.Slug)
			show = installedAppMap[item.Slug].Show
		}
		apps = append(apps, types.AppCenter{
			Name:        item.Name,
			Description: item.Description,
			Slug:        item.Slug,
			Channels: []struct {
				Slug      string `json:"slug"`
				Name      string `json:"name"`
				Panel     string `json:"panel"`
				Install   string `json:"-"`
				Uninstall string `json:"-"`
				Update    string `json:"-"`
				Subs      []struct {
					Log     string `json:"log"`
					Version string `json:"version"`
				} `json:"subs"`
			}(item.Channels),
			Installed:        installed,
			InstalledChannel: installedChannel,
			InstalledVersion: installedVersion,
			UpdateExist:      updateExist,
			Show:             show,
		})
	}

	paged, total := Paginate(r, apps)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *AppService) Install(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.App](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.appRepo.Install(req.Channel, req.Slug); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *AppService) Uninstall(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.AppSlug](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.appRepo.UnInstall(req.Slug); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *AppService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.AppSlug](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.appRepo.Update(req.Slug); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *AppService) UpdateShow(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.AppUpdateShow](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.appRepo.UpdateShow(req.Slug, req.Show); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *AppService) IsInstalled(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.AppSlug](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	app, err := s.appRepo.Get(req.Slug)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	installed, err := s.appRepo.IsInstalled(req.Slug)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, chix.M{
		"name":      app.Name,
		"installed": installed,
	})
}

func (s *AppService) UpdateCache(w http.ResponseWriter, r *http.Request) {
	if err := s.appRepo.UpdateCache(); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
