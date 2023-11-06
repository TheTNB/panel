package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UserShowAndDestroy struct {
	ID uint `form:"id" json:"id"`
}

func (r *UserShowAndDestroy) Authorize(ctx http.Context) error {
	return nil
}

func (r *UserShowAndDestroy) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "required|uint|min:1|exists:cert_users,id",
	}
}

func (r *UserShowAndDestroy) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UserShowAndDestroy) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UserShowAndDestroy) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
