package gitea

import (
	"fmt"
	"net/http"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read(fmt.Sprintf("%s/server/gitea/app.ini", app.Root))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, config)
}

func (s *Service) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = io.Write(fmt.Sprintf("%s/server/gitea/app.ini", app.Root), req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = systemctl.Restart("gitea"); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}
