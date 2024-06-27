package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Service struct {
	Service string `form:"service" json:"service"`
}

func (r *Service) Authorize(ctx http.Context) error {
	return nil
}

func (r *Service) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"service": "required|string|in:frps,frpc",
	}
}

func (r *Service) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Service) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Service) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Service) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
