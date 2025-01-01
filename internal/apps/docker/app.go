package docker

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/tnb-labs/panel/internal/service"
	"github.com/tnb-labs/panel/pkg/io"
	"github.com/tnb-labs/panel/pkg/systemctl"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (s *App) Route(r chi.Router) {
	r.Get("/config", s.GetConfig)
	r.Post("/config", s.UpdateConfig)
}

func (s *App) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read("/etc/docker/daemon.json")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, config)
}

func (s *App) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = io.Write("/etc/docker/daemon.json", req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = systemctl.Restart("docker"); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}
