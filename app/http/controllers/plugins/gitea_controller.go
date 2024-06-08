package plugins

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/app/http/controllers"
	requests "github.com/TheTNB/panel/app/http/requests/plugins/gitea"
	"github.com/TheTNB/panel/pkg/tools"
)

type GiteaController struct {
}

func NewGiteaController() *GiteaController {
	return &GiteaController{}
}

// Status
//
//	@Summary		服务状态
//	@Description	获取 Gitea 服务状态
//	@Tags			插件-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/status [get]
func (r *GiteaController) Status(ctx http.Context) http.Response {
	status, err := tools.ServiceStatus("gitea")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 Gitea 服务运行状态失败")
	}

	return controllers.Success(ctx, status)
}

// IsEnabled
//
//	@Summary		是否启用服务
//	@Description	获取是否启用 Gitea 服务
//	@Tags			插件-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/isEnabled [get]
func (r *GiteaController) IsEnabled(ctx http.Context) http.Response {
	enabled, err := tools.ServiceIsEnabled("gitea")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 Gitea 服务启用状态失败")
	}

	return controllers.Success(ctx, enabled)
}

// Enable
//
//	@Summary		启用服务
//	@Description	启用 Gitea 服务
//	@Tags			插件-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/enable [post]
func (r *GiteaController) Enable(ctx http.Context) http.Response {
	if err := tools.ServiceEnable("gitea"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "启用 Gitea 服务失败")
	}

	return controllers.Success(ctx, nil)
}

// Disable
//
//	@Summary		禁用服务
//	@Description	禁用 Gitea 服务
//	@Tags			插件-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/disable [post]
func (r *GiteaController) Disable(ctx http.Context) http.Response {
	if err := tools.ServiceDisable("gitea"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "禁用 Gitea 服务失败")
	}

	return controllers.Success(ctx, nil)
}

// Restart
//
//	@Summary		重启服务
//	@Description	重启 Gitea 服务
//	@Tags			插件-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/restart [post]
func (r *GiteaController) Restart(ctx http.Context) http.Response {
	if err := tools.ServiceRestart("gitea"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重启 Gitea 服务失败")
	}

	return controllers.Success(ctx, nil)
}

// Start
//
//	@Summary		启动服务
//	@Description	启动 Gitea 服务
//	@Tags			插件-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/start [post]
func (r *GiteaController) Start(ctx http.Context) http.Response {
	if err := tools.ServiceStart("gitea"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "启动 Gitea 服务失败")
	}

	status, err := tools.ServiceStatus("gitea")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 Gitea 服务运行状态失败")
	}

	return controllers.Success(ctx, status)
}

// Stop
//
//	@Summary		停止服务
//	@Description	停止 Gitea 服务
//	@Tags			插件-Gitea
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/gitea/stop [post]
func (r *GiteaController) Stop(ctx http.Context) http.Response {
	if err := tools.ServiceStop("gitea"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "停止 Gitea 服务失败")
	}

	status, err := tools.ServiceStatus("gitea")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 Gitea 服务运行状态失败")
	}

	return controllers.Success(ctx, !status)
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
	config, err := tools.Read("/www/server/gitea/app.ini")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
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
	sanitize := controllers.Sanitize(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.Write("/www/server/gitea/app.ini", updateRequest.Config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := tools.ServiceRestart("gitea"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}
