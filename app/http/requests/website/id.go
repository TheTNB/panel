package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ID struct {
	ID uint `form:"id" json:"id"`
}

func (r *ID) Authorize(ctx http.Context) error {
	return nil
}

func (r *ID) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "required|exists:websites,id",
	}
}

func (r *ID) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "uint",
	}
}

func (r *ID) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ID) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ID) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
