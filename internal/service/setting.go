package service

import "net/http"

type SettingService struct {
}

func NewSettingService() *SettingService {
	return &SettingService{}
}

func (s *SettingService) Get(w http.ResponseWriter, r *http.Request) {

}

func (s *SettingService) Update(w http.ResponseWriter, r *http.Request) {

}

func (s *SettingService) GetHttps(w http.ResponseWriter, r *http.Request) {

}

func (s *SettingService) UpdateHttps(w http.ResponseWriter, r *http.Request) {

}
