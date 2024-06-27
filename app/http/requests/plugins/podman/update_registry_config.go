package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateRegistryConfig struct {
	Config string `form:"config" json:"config"`
}

func (r *UpdateRegistryConfig) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateRegistryConfig) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"config": "required|string",
	}
}

func (r *UpdateRegistryConfig) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateRegistryConfig) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateRegistryConfig) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateRegistryConfig) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
