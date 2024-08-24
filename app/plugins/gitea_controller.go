package plugins

import (
	"github.com/goravel/framework/contracts/http"

	requests "github.com/TheTNB/panel/v2/app/http/requests/plugins/gitea"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type GiteaController struct {
}

func NewGiteaController() *GiteaController {
	return &GiteaController{}
}

// GetConfig
//
//	@Summary		获取配置
//	@Description	获取 Gitea 配置
//	@Tags			插件-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/config [get]
func (r *GiteaController) GetConfig(ctx http.Context) http.Response {
	config, err := io.Read("/www/server/gitea/app.ini")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
}

// UpdateConfig
//
//	@Summary		更新配置
//	@Description	更新 Gitea 配置
//	@Tags			插件-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.UpdateConfig	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/config [post]
func (r *GiteaController) UpdateConfig(ctx http.Context) http.Response {
	var updateRequest requests.UpdateConfig
	sanitize := h.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := io.Write("/www/server/gitea/app.ini", updateRequest.Config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := systemctl.Restart("gitea"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}
