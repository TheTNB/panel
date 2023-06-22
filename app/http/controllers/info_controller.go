package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/models"
)

type InfoController struct {
	//Dependent services
}

func NewInfoController() *InfoController {
	return &InfoController{
		//Inject services
	}
}

func (r *InfoController) Name(ctx http.Context) {
	var setting models.Setting
	err := facades.Orm().Query().Where("key", "name").First(&setting)
	if err != nil {
		facades.Log().Error("[面板][InfoController] 查询面板名称失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, http.Json{
		"name": setting.Value,
	})
}
