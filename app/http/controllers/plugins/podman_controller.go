package plugins

import (
	"github.com/goravel/framework/contracts/http"

	requests "github.com/TheTNB/panel/v2/app/http/requests/plugins/podman"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type PodmanController struct {
}

func NewPodmanController() *PodmanController {
	return &PodmanController{}
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
	config, err := io.Read("/etc/containers/registries.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
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
	sanitize := h.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := io.Write("/etc/containers/registries.conf", updateRequest.Config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := systemctl.Restart("podman"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
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
	config, err := io.Read("/etc/containers/storage.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
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
	sanitize := h.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := io.Write("/etc/containers/storage.conf", updateRequest.Config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := systemctl.Restart("podman"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}
