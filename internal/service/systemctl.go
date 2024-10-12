package service

import (
	"net/http"

	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type SystemctlService struct {
}

func NewSystemctlService() *SystemctlService {
	return &SystemctlService{}
}

func (s *SystemctlService) Status(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemctlService](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	status, err := systemctl.Status(req.Service)
	if err != nil {
		Error(w, http.StatusInternalServerError, "获取 %s 服务运行状态失败", req.Service)
		return
	}

	Success(w, status)
}

func (s *SystemctlService) IsEnabled(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemctlService](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	enabled, err := systemctl.IsEnabled(req.Service)
	if err != nil {
		Error(w, http.StatusInternalServerError, "获取 %s 服务启用状态失败", req.Service)
		return
	}

	Success(w, enabled)
}

func (s *SystemctlService) Enable(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemctlService](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = systemctl.Enable(req.Service); err != nil {
		Error(w, http.StatusInternalServerError, "启用 %s 服务失败", req.Service)
		return
	}

	Success(w, nil)
}

func (s *SystemctlService) Disable(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemctlService](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = systemctl.Disable(req.Service); err != nil {
		Error(w, http.StatusInternalServerError, "禁用 %s 服务失败", req.Service)
		return
	}

	Success(w, nil)
}

func (s *SystemctlService) Restart(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemctlService](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = systemctl.Restart(req.Service); err != nil {
		Error(w, http.StatusInternalServerError, "重启 %s 服务失败", req.Service)
		return
	}

	Success(w, nil)
}

func (s *SystemctlService) Reload(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemctlService](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = systemctl.Reload(req.Service); err != nil {
		Error(w, http.StatusInternalServerError, "重载 %s 服务失败", req.Service)
		return
	}

	Success(w, nil)
}

func (s *SystemctlService) Start(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemctlService](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = systemctl.Start(req.Service); err != nil {
		Error(w, http.StatusInternalServerError, "启动 %s 服务失败", req.Service)
		return
	}

	Success(w, nil)
}

func (s *SystemctlService) Stop(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemctlService](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = systemctl.Stop(req.Service); err != nil {
		Error(w, http.StatusInternalServerError, "停止 %s 服务失败", req.Service)
		return
	}

	Success(w, nil)
}
