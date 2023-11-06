package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Renew struct {
	ID uint `form:"id" json:"id"`
}

func (r *Renew) Authorize(ctx http.Context) error {
	return nil
}

func (r *Renew) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "required|exists:certs,id",
	}
}

func (r *Renew) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Renew) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Renew) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
