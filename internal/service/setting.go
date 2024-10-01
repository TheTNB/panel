package service

import (
	"net/http"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/tools"
)

type SettingService struct {
	settingRepo biz.SettingRepo
}

func NewSettingService() *SettingService {
	return &SettingService{
		settingRepo: data.NewSettingRepo(),
	}
}

func (s *SettingService) Get(w http.ResponseWriter, r *http.Request) {
	setting, err := s.settingRepo.GetPanelSetting(r.Context())
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, setting)
}

func (s *SettingService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.PanelSetting](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	restart := false
	if restart, err = s.settingRepo.UpdatePanelSetting(r.Context(), req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if restart {
		tools.RestartPanel()
	}

	Success(w, nil)
}
