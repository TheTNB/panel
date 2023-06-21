package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type OpenRestyController struct {
	//Dependent services
}

func NewOpenrestyController() *OpenRestyController {
	return &OpenRestyController{
		//Inject services
	}
}

func (r *OpenRestyController) Show(ctx http.Context) {
	ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
	})
}
