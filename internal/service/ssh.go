package service

import (
	"net/http"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
)

type SSHService struct {
	sshRepo biz.SSHRepo
}

func NewSSHService() *SSHService {
	return &SSHService{
		sshRepo: data.NewSSHRepo(),
	}
}

func (s *SSHService) GetInfo(w http.ResponseWriter, r *http.Request) {
	info, err := s.sshRepo.GetInfo()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, info)
}

func (s *SSHService) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SSHUpdateInfo](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.sshRepo.UpdateInfo(req); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
}
