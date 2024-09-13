package service

import "net/http"

type SSHService struct {
}

func NewSSHService() *SSHService {
	return &SSHService{}
}

func (s *SSHService) GetInfo(w http.ResponseWriter, r *http.Request) {

}

func (s *SSHService) UpdateInfo(w http.ResponseWriter, r *http.Request) {

}

func (s *SSHService) Session(w http.ResponseWriter, r *http.Request) {

}
