package controllers

import (
	"github.com/gofiber/swagger"
	"github.com/goravel/fiber"
	"github.com/goravel/framework/contracts/http"

	_ "panel/docs"
)

type SwaggerController struct {
	// Dependent services
}

// Config stores fiberSwagger configuration variables.
type Config struct {
	URL                  string
	InstanceName         string
	DocExpansion         string
	DomID                string
	DeepLinking          bool
	PersistAuthorization bool
}

func NewSwaggerController() *SwaggerController {
	return &SwaggerController{
		// Inject services
	}
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
	err := swagger.New(swagger.Config{
		Title: "耗子面板 Swagger",
	})(ctx.(*fiber.Context).Instance())
	if err != nil {
		return Error(ctx, http.StatusNotFound, err.Error())
	}

	return nil
}
