package service

import (
	"net/http"
	"path/filepath"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/pkg/io"
)

type WebsiteService struct {
	websiteRepo biz.WebsiteRepo
	settingRepo biz.SettingRepo
}

func NewWebsiteService() *WebsiteService {
	return &WebsiteService{
		websiteRepo: data.NewWebsiteRepo(),
		settingRepo: data.NewSettingRepo(),
	}
}

func (s *WebsiteService) GetDefaultConfig(w http.ResponseWriter, r *http.Request) {
	index, err := io.Read(filepath.Join(panel.Root, "server/openresty/html/index.html"))
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	stop, err := io.Read(filepath.Join(panel.Root, "server/openresty/html/stop.html"))
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, chix.M{
		"index": index,
		"stop":  stop,
	})
}

func (s *WebsiteService) UpdateDefaultConfig(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteDefaultConfig](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.UpdateDefaultConfig(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *WebsiteService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	websites, total, err := s.websiteRepo.List(req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": websites,
	})
}

func (s *WebsiteService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if len(req.Path) == 0 {
		req.Path, _ = s.settingRepo.Get(biz.SettingKeyWebsitePath)
		req.Path = filepath.Join(req.Path, req.Name)
	}

	if _, err = s.websiteRepo.Create(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *WebsiteService) Get(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	config, err := s.websiteRepo.Get(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, config)
}

func (s *WebsiteService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteUpdate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.Update(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *WebsiteService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteDelete](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.Delete(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *WebsiteService) ClearLog(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.ClearLog(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *WebsiteService) UpdateRemark(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteUpdateRemark](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.UpdateRemark(req.ID, req.Remark); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *WebsiteService) ResetConfig(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.ResetConfig(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *WebsiteService) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteUpdateStatus](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.UpdateStatus(req.ID, req.Status); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
