package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateConfig struct {
	Config string `form:"config" json:"config"`
}

func (r *UpdateConfig) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateConfig) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"config": "required|string",
	}
}

func (r *UpdateConfig) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateConfig) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateConfig) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
