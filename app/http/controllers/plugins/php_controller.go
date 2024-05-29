package plugins

import (
	"github.com/TheTNB/panel/app/http/controllers"
	"github.com/TheTNB/panel/internal"
	"github.com/TheTNB/panel/internal/services"
	"github.com/goravel/framework/contracts/http"
)

type PHPController struct {
	service internal.PHP
}

func NewPHPController(version uint) *PHPController {
	return &PHPController{
		service: services.NewPHPImpl(version),
	}
}

func (r *PHPController) Status(ctx http.Context) http.Response {
	status, err := r.service.Status()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, status)
}

func (r *PHPController) Reload(ctx http.Context) http.Response {
	if err := r.service.Reload(); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) Start(ctx http.Context) http.Response {
	if err := r.service.Start(); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) Stop(ctx http.Context) http.Response {
	if err := r.service.Stop(); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) Restart(ctx http.Context) http.Response {
	if err := r.service.Restart(); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) GetConfig(ctx http.Context) http.Response {
	config, err := r.service.GetConfig()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

func (r *PHPController) SaveConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	if err := r.service.SaveConfig(config); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) GetFPMConfig(ctx http.Context) http.Response {
	config, err := r.service.GetFPMConfig()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

func (r *PHPController) SaveFPMConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	if err := r.service.SaveFPMConfig(config); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) Load(ctx http.Context) http.Response {
	load, err := r.service.Load()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, load)
}

func (r *PHPController) ErrorLog(ctx http.Context) http.Response {
	log, _ := r.service.GetErrorLog()
	return controllers.Success(ctx, log)
}

func (r *PHPController) SlowLog(ctx http.Context) http.Response {
	log, _ := r.service.GetSlowLog()
	return controllers.Success(ctx, log)
}

func (r *PHPController) ClearErrorLog(ctx http.Context) http.Response {
	err := r.service.ClearErrorLog()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) ClearSlowLog(ctx http.Context) http.Response {
	err := r.service.ClearSlowLog()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) GetExtensionList(ctx http.Context) http.Response {
	extensions, err := r.service.GetExtensions()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, extensions)
}

func (r *PHPController) InstallExtension(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	if err := r.service.InstallExtension(slug); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

func (r *PHPController) UninstallExtension(ctx http.Context) http.Response {
	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	if err := r.service.UninstallExtension(slug); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}
