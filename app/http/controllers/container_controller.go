package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type ContainerController struct {
	// Dependent services
}

func NewContainerController() *ContainerController {
	return &ContainerController{
		// Inject services
	}
}

func (r *ContainerController) Index(ctx http.Context) http.Response {
	return nil
}
