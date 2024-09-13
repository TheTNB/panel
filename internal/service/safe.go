package service

import "net/http"

type SafeService struct {
}

func NewSafeService() *SafeService {
	return &SafeService{}
}

func (s *SafeService) GetSSHStatus(w http.ResponseWriter, r *http.Request) {

}

func (s *SafeService) UpdateSSHStatus(w http.ResponseWriter, r *http.Request) {

}

func (s *SafeService) GetPingStatus(w http.ResponseWriter, r *http.Request) {

}

func (s *SafeService) UpdatePingStatus(w http.ResponseWriter, r *http.Request) {

}
