package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/swaggo/http-swagger/v2"

	_ "github.com/TheTNB/panel/docs"
)

type SwaggerController struct {
	// Dependent services
}

func NewSwaggerController() *SwaggerController {
	return &SwaggerController{}
}

// Index
//
//	@Summary		Swagger UI
//	@Description	Swagger UI
//	@Tags			Swagger
//	@Success		200
//	@Failure		500
//	@Router			/swagger [get]
func (r *SwaggerController) Index(ctx http.Context) http.Response {
	if !facades.Config().GetBool("app.debug") {
		return Error(ctx, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	handler := httpSwagger.Handler()
	handler(ctx.Response().Writer(), ctx.Request().Origin())

	return nil
}
