package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type SystemController struct {
}

func NewSystemController() *SystemController {
	return &SystemController{}
}

// ServiceStatus
//
//	@Summary	服务状态
//	@Tags		系统
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	query		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/system/service/status [get]
func (r *SystemController) ServiceStatus(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"service": "required|string",
	}); sanitize != nil {
		return sanitize
	}

	service := ctx.Request().Query("service")
	status, err := systemctl.Status(service)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("获取 %s 服务运行状态失败", service))
	}

	return h.Success(ctx, status)
}

// ServiceIsEnabled
//
//	@Summary	是否启用服务
//	@Tags		系统
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	query		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/system/service/isEnabled [get]
func (r *SystemController) ServiceIsEnabled(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"service": "required|string",
	}); sanitize != nil {
		return sanitize
	}

	service := ctx.Request().Query("service")
	enabled, err := systemctl.IsEnabled(service)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("获取 %s 服务启用状态失败", service))
	}

	return h.Success(ctx, enabled)
}

// ServiceEnable
//
//	@Summary	启用服务
//	@Tags		系统
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/system/service/enable [post]
func (r *SystemController) ServiceEnable(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"service": "required|string",
	}); sanitize != nil {
		return sanitize
	}

	service := ctx.Request().Input("service")
	if err := systemctl.Enable(service); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("启用 %s 服务失败", service))
	}

	return h.Success(ctx, nil)
}

// ServiceDisable
//
//	@Summary	禁用服务
//	@Tags		系统
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/system/service/disable [post]
func (r *SystemController) ServiceDisable(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"service": "required|string",
	}); sanitize != nil {
		return sanitize
	}

	service := ctx.Request().Input("service")
	if err := systemctl.Disable(service); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("禁用 %s 服务失败", service))
	}

	return h.Success(ctx, nil)
}

// ServiceRestart
//
//	@Summary	重启服务
//	@Tags		系统
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/system/service/restart [post]
func (r *SystemController) ServiceRestart(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"service": "required|string",
	}); sanitize != nil {
		return sanitize
	}

	service := ctx.Request().Input("service")
	if err := systemctl.Restart(service); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重启 %s 服务失败", service))
	}

	return h.Success(ctx, nil)
}

// ServiceReload
//
//	@Summary	重载服务
//	@Tags		系统
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/system/service/reload [post]
func (r *SystemController) ServiceReload(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"service": "required|string",
	}); sanitize != nil {
		return sanitize
	}

	service := ctx.Request().Input("service")
	if err := systemctl.Reload(service); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重载 %s 服务失败", service))
	}

	return h.Success(ctx, nil)
}

// ServiceStart
//
//	@Summary	启动服务
//	@Tags		系统
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/system/service/start [post]
func (r *SystemController) ServiceStart(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"service": "required|string",
	}); sanitize != nil {
		return sanitize
	}

	service := ctx.Request().Input("service")
	if err := systemctl.Start(service); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("启动 %s 服务失败", service))
	}

	return h.Success(ctx, nil)
}

// ServiceStop
//
//	@Summary	停止服务
//	@Tags		系统
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		string	true	"request"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/panel/system/service/stop [post]
func (r *SystemController) ServiceStop(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"service": "required|string",
	}); sanitize != nil {
		return sanitize
	}

	service := ctx.Request().Input("service")
	if err := systemctl.Stop(service); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("停止 %s 服务失败", service))
	}

	return h.Success(ctx, nil)
}
