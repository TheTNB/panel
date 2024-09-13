package request

import (
	"net/http"
)

type Paginate struct {
	Page  uint `json:"page" form:"page" query:"page" validate:"required,number,gte=1"`
	Limit uint `json:"limit" form:"limit" query:"limit" validate:"required,number,gte=1,lte=1000"`
}

func (r *Paginate) Messages(_ *http.Request) map[string]string {
	return map[string]string{
		"Page.gte":       "页码必须大于或等于1",
		"Limit.gte":      "每页数量必须大于或等于1",
		"Limit.lte":      "每页数量必须小于或等于1000",
		"Page.number":    "页码必须是数字",
		"Limit.number":   "每页数量必须是数字",
		"Page.required":  "页码不能为空",
		"Limit.required": "每页数量不能为空",
	}
}

func (r *Paginate) Prepare(_ *http.Request) error {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	return nil
}
