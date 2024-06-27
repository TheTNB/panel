package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type DNSShowAndDestroy struct {
	ID uint `form:"id" json:"id"`
}

func (r *DNSShowAndDestroy) Authorize(ctx http.Context) error {
	return nil
}

func (r *DNSShowAndDestroy) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "required|uint|min:1|exists:cert_dns,id",
	}
}

func (r *DNSShowAndDestroy) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSShowAndDestroy) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSShowAndDestroy) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSShowAndDestroy) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
