package plugins

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/app/http/controllers"
	requests "github.com/TheTNB/panel/app/http/requests/plugins/frp"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type FrpController struct {
}

func NewFrpController() *FrpController {
	return &FrpController{}
}

// GetConfig
//
//	@Summary		获取配置
//	@Description	获取 Frp 配置
//	@Tags			插件-Frp
//	@Produce		json
//	@Security		BearerToken
//	@Param			service	query		string	false	"服务"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/frp/config [get]
func (r *FrpController) GetConfig(ctx http.Context) http.Response {
	var serviceRequest requests.Service
	sanitize := controllers.SanitizeRequest(ctx, &serviceRequest)
	if sanitize != nil {
		return sanitize
	}

	config, err := io.Read(fmt.Sprintf("/www/server/frp/%s.toml", serviceRequest.Service))
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

// UpdateConfig
//
//	@Summary		更新配置
//	@Description	更新 Frp 配置
//	@Tags			插件-Frp
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.UpdateConfig	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/frp/config [post]
func (r *FrpController) UpdateConfig(ctx http.Context) http.Response {
	var updateRequest requests.UpdateConfig
	sanitize := controllers.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := io.Write(fmt.Sprintf("/www/server/frp/%s.toml", updateRequest.Service), updateRequest.Config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := systemctl.Restart(updateRequest.Service); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}
