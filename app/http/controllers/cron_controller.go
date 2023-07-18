package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type CronController struct {
	//Dependent services
}

func NewCronController() *CronController {
	return &CronController{
		//Inject services
	}
}

func (r *CronController) Index(ctx http.Context) {
}
