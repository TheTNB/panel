package plugins

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/app/http/controllers"
	requests "github.com/TheTNB/panel/app/http/requests/plugins/podman"
	"github.com/TheTNB/panel/pkg/tools"
)

type PodmanController struct {
}

func NewPodmanController() *PodmanController {
	return &PodmanController{}
}

// Status
//
//	@Summary		服务状态
//	@Description	获取 Podman 服务状态
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/status [get]
func (r *PodmanController) Status(ctx http.Context) http.Response {
	status, err := tools.ServiceStatus("podman")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 Podman 服务运行状态失败")
	}

	return controllers.Success(ctx, status)
}

// IsEnabled
//
//	@Summary		是否启用服务
//	@Description	获取是否启用 Podman 服务
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/isEnabled [get]
func (r *PodmanController) IsEnabled(ctx http.Context) http.Response {
	enabled, err := tools.ServiceIsEnabled("podman")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 Podman 服务启用状态失败")
	}

	return controllers.Success(ctx, enabled)
}

// Enable
//
//	@Summary		启用服务
//	@Description	启用 Podman 服务
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/enable [post]
func (r *PodmanController) Enable(ctx http.Context) http.Response {
	if err := tools.ServiceEnable("podman"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "启用 Podman 服务失败")
	}

	return controllers.Success(ctx, nil)
}

// Disable
//
//	@Summary		禁用服务
//	@Description	禁用 Podman 服务
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/disable [post]
func (r *PodmanController) Disable(ctx http.Context) http.Response {
	if err := tools.ServiceDisable("podman"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "禁用 Podman 服务失败")
	}

	return controllers.Success(ctx, nil)
}

// Restart
//
//	@Summary		重启服务
//	@Description	重启 Podman 服务
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/restart [post]
func (r *PodmanController) Restart(ctx http.Context) http.Response {
	if err := tools.ServiceRestart("podman"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重启 Podman 服务失败")
	}

	return controllers.Success(ctx, nil)
}

// Start
//
//	@Summary		启动服务
//	@Description	启动 Podman 服务
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/start [post]
func (r *PodmanController) Start(ctx http.Context) http.Response {
	if err := tools.ServiceStart("podman"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "启动 Podman 服务失败")
	}

	status, err := tools.ServiceStatus("podman")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 Podman 服务运行状态失败")
	}

	return controllers.Success(ctx, status)
}

// Stop
//
//	@Summary		停止服务
//	@Description	停止 Podman 服务
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/stop [post]
func (r *PodmanController) Stop(ctx http.Context) http.Response {
	if err := tools.ServiceStop("podman"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "停止 Podman 服务失败")
	}

	status, err := tools.ServiceStatus("podman")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 Podman 服务运行状态失败")
	}

	return controllers.Success(ctx, !status)
}

// GetRegistryConfig
//
//	@Summary		获取注册表配置
//	@Description	获取 Podman 注册表配置
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/registryConfig [get]
func (r *PodmanController) GetRegistryConfig(ctx http.Context) http.Response {
	config, err := tools.Read("/etc/containers/registries.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

// UpdateRegistryConfig
//
//	@Summary		更新注册表配置
//	@Description	更新 Podman 注册表配置
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.UpdateRegistryConfig	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/registryConfig [post]
func (r *PodmanController) UpdateRegistryConfig(ctx http.Context) http.Response {
	var updateRequest requests.UpdateRegistryConfig
	sanitize := controllers.Sanitize(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.Write("/etc/containers/registries.conf", updateRequest.Config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := tools.ServiceRestart("podman"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

// GetStorageConfig
//
//	@Summary		获取存储配置
//	@Description	获取 Podman 存储配置
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/storageConfig [get]
func (r *PodmanController) GetStorageConfig(ctx http.Context) http.Response {
	config, err := tools.Read("/etc/containers/storage.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

// UpdateStorageConfig
//
//	@Summary		更新存储配置
//	@Description	更新 Podman 存储配置
//	@Tags			插件-Podman
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.UpdateStorageConfig	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/podman/storageConfig [post]
func (r *PodmanController) UpdateStorageConfig(ctx http.Context) http.Response {
	var updateRequest requests.UpdateStorageConfig
	sanitize := controllers.Sanitize(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.Write("/etc/containers/storage.conf", updateRequest.Config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := tools.ServiceRestart("podman"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}
