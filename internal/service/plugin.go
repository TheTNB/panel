package service

import "net/http"

type PluginService struct {
}

func NewPluginService() *PluginService {
	return &PluginService{}
}

func (s *PluginService) List(w http.ResponseWriter, r *http.Request) {

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
