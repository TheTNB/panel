package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type CertShowAndDestroy struct {
	ID uint `form:"id" json:"id"`
}

func (r *CertShowAndDestroy) Authorize(ctx http.Context) error {
	return nil
}

func (r *CertShowAndDestroy) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "required|uint|min:1|exists:certs,id",
	}
}

func (r *CertShowAndDestroy) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertShowAndDestroy) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertShowAndDestroy) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
