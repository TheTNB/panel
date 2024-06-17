package plugins

import (
	"github.com/TheTNB/panel/app/http/controllers"
	"github.com/TheTNB/panel/internal/services"
	"github.com/goravel/framework/contracts/http"
)

type PHPController struct{}

func NewPHPController() *PHPController {
	return &PHPController{}
}

func (r *PHPController) GetConfig(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	config, err := service.GetConfig()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

func (r *PHPController) SaveConfig(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	config := ctx.Request().Input("config")
	if err := service.SaveConfig(config); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) GetFPMConfig(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	config, err := service.GetFPMConfig()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

func (r *PHPController) SaveFPMConfig(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	config := ctx.Request().Input("config")
	if err := service.SaveFPMConfig(config); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) Load(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	load, err := service.Load()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, load)
}

func (r *PHPController) ErrorLog(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	log, _ := service.GetErrorLog()
	return controllers.Success(ctx, log)
}

func (r *PHPController) SlowLog(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	log, _ := service.GetSlowLog()
	return controllers.Success(ctx, log)
}

func (r *PHPController) ClearErrorLog(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	err := service.ClearErrorLog()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) ClearSlowLog(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	err := service.ClearSlowLog()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) GetExtensionList(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	extensions, err := service.GetExtensions()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, extensions)
}

func (r *PHPController) InstallExtension(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	if err := service.InstallExtension(slug); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) UninstallExtension(ctx http.Context) http.Response {
	service := services.NewPHPImpl(uint(ctx.Request().RouteInt("version")))
	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	if err := service.UninstallExtension(slug); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}
