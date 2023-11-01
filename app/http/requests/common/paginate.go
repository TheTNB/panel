package commonrequests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Paginate struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}

func (r *Paginate) Authorize(ctx http.Context) error {
	return nil
}

func (r *Paginate) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"page":  "required|uint|min:1",
		"limit": "required|uint|min:1",
	}
}

func (r *Paginate) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"page.required":  "分页参数 page 不能为空",
		"page.uint":      "分页参数 page 必须是一个整数",
		"page.min":       "分页参数 page 必须大于等于 1",
		"limit.required": "分页参数 limit 不能为空",
		"limit.uint":     "分页参数 limit 必须是一个整数",
		"limit.min":      "分页参数 limit 必须大于等于 1",
	}
}

func (r *Paginate) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Paginate) PrepareForValidation(ctx http.Context, data validation.Data) error {
	_ = data.Set("page", ctx.Request().QueryInt("page"))
	_ = data.Set("limit", ctx.Request().QueryInt("limit"))

	return nil
}
