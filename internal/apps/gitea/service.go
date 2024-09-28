package gitea

import (
	"fmt"
	"net/http"

	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// GetConfig
//
//	@Summary		获取配置
//	@Tags			应用-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/config [get]
func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read(fmt.Sprintf("%s/server/gitea/app.ini", panel.Root))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, config)
}

// UpdateConfig
//
//	@Summary		更新配置
//	@Tags			应用-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.UpdateConfig	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/config [post]
func (s *Service) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = io.Write(fmt.Sprintf("%s/server/gitea/app.ini", panel.Root), req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = systemctl.Restart("gitea"); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}
