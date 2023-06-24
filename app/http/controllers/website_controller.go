package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"panel/app/models"
)

type WebsiteController struct {
	//Dependent services
}

func NewWebsiteController() *WebsiteController {
	return &WebsiteController{
		//Inject services
	}
}

func (r *WebsiteController) List(ctx http.Context) {
	limit := ctx.Request().QueryInt("limit")
	page := ctx.Request().QueryInt("page")

	var websites []models.Website
	var total int64
	err := facades.Orm().Query().Paginate(page, limit, &websites, &total)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 查询网站列表失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, http.Json{
		"total": total,
		"items": websites,
	})
}
