package service

import "net/http"

type SystemctlService struct {
}

func NewSystemctlService() *SystemctlService {
	return &SystemctlService{}
}

func (s *SystemctlService) Status(w http.ResponseWriter, r *http.Request) {

}

func (s *SystemctlService) IsEnabled(w http.ResponseWriter, r *http.Request) {

}

func (s *SystemctlService) Enable(w http.ResponseWriter, r *http.Request) {

}

func (s *SystemctlService) Disable(w http.ResponseWriter, r *http.Request) {

}

func (s *SystemctlService) Restart(w http.ResponseWriter, r *http.Request) {

}

func (s *SystemctlService) Reload(w http.ResponseWriter, r *http.Request) {

}

func (s *SystemctlService) Start(w http.ResponseWriter, r *http.Request) {

}

func (s *SystemctlService) Stop(w http.ResponseWriter, r *http.Request) {

}
