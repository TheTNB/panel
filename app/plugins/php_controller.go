package plugins

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
)

type PHPController struct{}

func NewPHPController() *PHPController {
	return &PHPController{}
}

// GetConfig
//
//	@Summary	获取配置
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/config [get]
func (r *PHPController) GetConfig(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	config, err := service.GetConfig()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
}

// SaveConfig
//
//	@Summary	保存配置
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int		true	"PHP 版本"
//	@Param		config	body		string	true	"配置"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/config [post]
func (r *PHPController) SaveConfig(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	config := ctx.Request().Input("config")
	if err := service.SaveConfig(config); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// GetFPMConfig
//
//	@Summary	获取 FPM 配置
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/fpmConfig [get]
func (r *PHPController) GetFPMConfig(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	config, err := service.GetFPMConfig()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
}

// SaveFPMConfig
//
//	@Summary	保存 FPM 配置
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int		true	"PHP 版本"
//	@Param		config	body		string	true	"配置"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/fpmConfig [post]
func (r *PHPController) SaveFPMConfig(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	config := ctx.Request().Input("config")
	if err := service.SaveFPMConfig(config); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// Load
//
//	@Summary	获取负载状态
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/load [get]
func (r *PHPController) Load(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	load, err := service.Load()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, load)
}

// ErrorLog
//
//	@Summary	获取错误日志
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/errorLog [get]
func (r *PHPController) ErrorLog(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	log, _ := service.GetErrorLog()
	return h.Success(ctx, log)
}

// SlowLog
//
//	@Summary	获取慢日志
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/slowLog [get]
func (r *PHPController) SlowLog(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	log, _ := service.GetSlowLog()
	return h.Success(ctx, log)
}

// ClearErrorLog
//
//	@Summary	清空错误日志
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/clearErrorLog [post]
func (r *PHPController) ClearErrorLog(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	err := service.ClearErrorLog()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ClearSlowLog
//
//	@Summary	清空慢日志
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/clearSlowLog [post]
func (r *PHPController) ClearSlowLog(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	err := service.ClearSlowLog()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ExtensionList
//
//	@Summary	获取扩展列表
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/extensions [get]
func (r *PHPController) ExtensionList(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	extensions, err := service.GetExtensions()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, extensions)
}

// InstallExtension
//
//	@Summary	安装扩展
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int		true	"PHP 版本"
//	@Param		slug	query		string	true	"slug"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/extensions [post]
func (r *PHPController) InstallExtension(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	if err := service.InstallExtension(slug); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// UninstallExtension
//
//	@Summary	卸载扩展
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int		true	"PHP 版本"
//	@Param		slug	query		string	true	"slug"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/extensions [delete]
func (r *PHPController) UninstallExtension(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	if err := service.UninstallExtension(slug); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}
