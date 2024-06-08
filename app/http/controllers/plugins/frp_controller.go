package plugins

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/app/http/controllers"
	requests "github.com/TheTNB/panel/app/http/requests/plugins/frp"
	"github.com/TheTNB/panel/pkg/tools"
)

type FrpController struct {
}

func NewFrpController() *FrpController {
	return &FrpController{}
}

// Status
//
//	@Summary		服务状态
//	@Description	获取 Frp 服务状态
//	@Tags			插件-Frp
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/frp/status [get]
func (r *FrpController) Status(ctx http.Context) http.Response {
	frps, err := tools.ServiceStatus("frps")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 frps 服务运行状态失败")
	}
	frpc, err := tools.ServiceStatus("frpc")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 frpc 服务运行状态失败")
	}

	return controllers.Success(ctx, http.Json{
		"frps": frps,
		"frpc": frpc,
	})
}

// IsEnabled
//
//	@Summary		是否启用服务
//	@Description	获取是否启用 Frp 服务
//	@Tags			插件-Frp
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/frp/isEnabled [get]
func (r *FrpController) IsEnabled(ctx http.Context) http.Response {
	frps, err := tools.ServiceIsEnabled("frps")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 frps 服务启用状态失败")
	}
	frpc, err := tools.ServiceIsEnabled("frpc")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 frpc 服务启用状态失败")
	}

	return controllers.Success(ctx, http.Json{
		"frps": frps,
		"frpc": frpc,
	})
}

// Enable
//
//	@Summary		启用服务
//	@Description	启用 Frp 服务
//	@Tags			插件-Frp
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Service	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/frp/enable [post]
func (r *FrpController) Enable(ctx http.Context) http.Response {
	var serviceRequest requests.Service
	sanitize := controllers.Sanitize(ctx, &serviceRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.ServiceEnable(serviceRequest.Service); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("启用 %s 服务失败", serviceRequest.Service))
	}

	return controllers.Success(ctx, nil)
}

// Disable
//
//	@Summary		禁用服务
//	@Description	禁用 Frp 服务
//	@Tags			插件-Frp
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Service	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/frp/disable [post]
func (r *FrpController) Disable(ctx http.Context) http.Response {
	var serviceRequest requests.Service
	sanitize := controllers.Sanitize(ctx, &serviceRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.ServiceDisable(serviceRequest.Service); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("禁用 %s 服务失败", serviceRequest.Service))
	}

	return controllers.Success(ctx, nil)
}

// Restart
//
//	@Summary		重启服务
//	@Description	重启 Frp 服务
//	@Tags			插件-Frp
//	@Produce		json
//	@Security		BearerToken
//	@param			data	body		requests.Service	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/frp/restart [post]
func (r *FrpController) Restart(ctx http.Context) http.Response {
	var serviceRequest requests.Service
	sanitize := controllers.Sanitize(ctx, &serviceRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.ServiceRestart(serviceRequest.Service); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重启 %s 服务失败", serviceRequest.Service))
	}

	return controllers.Success(ctx, nil)
}

// Start
//
//	@Summary		启动服务
//	@Description	启动 Frp 服务
//	@Tags			插件-Frp
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Service	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/frp/start [post]
func (r *FrpController) Start(ctx http.Context) http.Response {
	var serviceRequest requests.Service
	sanitize := controllers.Sanitize(ctx, &serviceRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.ServiceStart(serviceRequest.Service); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("启动 %s 服务失败", serviceRequest.Service))
	}

	status, err := tools.ServiceStatus(serviceRequest.Service)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("获取 %s 服务运行状态失败", serviceRequest.Service))
	}

	return controllers.Success(ctx, status)
}

// Stop
//
//	@Summary		停止服务
//	@Description	停止 Frp 服务
//	@Tags			插件-Frp
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Service	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/frp/stop [post]
func (r *FrpController) Stop(ctx http.Context) http.Response {
	var serviceRequest requests.Service
	sanitize := controllers.Sanitize(ctx, &serviceRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.ServiceStop(serviceRequest.Service); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("停止 %s 服务失败", serviceRequest.Service))
	}

	status, err := tools.ServiceStatus(serviceRequest.Service)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("获取 %s 服务运行状态失败", serviceRequest.Service))
	}

	return controllers.Success(ctx, !status)
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
	sanitize := controllers.Sanitize(ctx, &serviceRequest)
	if sanitize != nil {
		return sanitize
	}

	config, err := tools.Read(fmt.Sprintf("/www/server/frp/%s.toml", serviceRequest.Service))
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
	sanitize := controllers.Sanitize(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.Write(fmt.Sprintf("/www/server/frp/%s.toml", updateRequest.Service), updateRequest.Config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := tools.ServiceRestart(updateRequest.Service); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}
