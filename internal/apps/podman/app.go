package podman

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
	r.Get("/registryConfig", s.GetRegistryConfig)
	r.Post("/registryConfig", s.UpdateRegistryConfig)
	r.Get("/storageConfig", s.GetStorageConfig)
	r.Post("/storageConfig", s.UpdateStorageConfig)
}

func (s *App) GetRegistryConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read("/etc/containers/registries.conf")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, config)
}

func (s *App) UpdateRegistryConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = io.Write("/etc/containers/registries.conf", req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = systemctl.Restart("podman"); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}

func (s *App) GetStorageConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read("/etc/containers/storage.conf")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, config)
}

func (s *App) UpdateStorageConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = io.Write("/etc/containers/storage.conf", req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = systemctl.Restart("podman"); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}
