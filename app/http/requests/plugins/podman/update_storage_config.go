package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateStorageConfig struct {
	Config string `form:"config" json:"config"`
}

func (r *UpdateStorageConfig) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateStorageConfig) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"config": "required|string",
	}
}

func (r *UpdateStorageConfig) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateStorageConfig) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateStorageConfig) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateStorageConfig) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
