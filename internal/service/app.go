package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/str"
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
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	installedAppMap := make(map[string]*biz.App)

	for _, p := range installedApps {
		installedAppMap[p.Slug] = p
	}

	var apps []types.StoreApp
	for _, item := range all {
		installed, installedVersion, currentVersion, show := false, "", "", false
		if str.FirstElement(item.Versions) != nil {
			currentVersion = str.FirstElement(item.Versions).Version
		}
		if _, ok := installedAppMap[item.Slug]; ok {
			installed = true
			installedVersion = installedAppMap[item.Slug].Version
			show = installedAppMap[item.Slug].Show
		}
		apps = append(apps, types.StoreApp{
			Name:             item.Name,
			Description:      item.Description,
			Slug:             item.Slug,
			Version:          currentVersion,
			Installed:        installed,
			InstalledVersion: installedVersion,
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
	req, err := Bind[request.AppSlug](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.appRepo.Install(req.Slug); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
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

	if err = s.appRepo.Uninstall(req.Slug); err != nil {
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
