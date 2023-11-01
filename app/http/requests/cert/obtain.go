package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Obtain struct {
	ID uint `form:"id" json:"id"`
}

func (r *Obtain) Authorize(ctx http.Context) error {
	return nil
}

func (r *Obtain) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "required|exists:certs,id",
	}
}

func (r *Obtain) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"id.required": "证书 ID 不能为空",
		"id.exists":   "证书 ID 不存在",
	}
}

func (r *Obtain) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Obtain) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
