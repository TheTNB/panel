package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ID struct {
	ID string `form:"id" json:"id"`
}

func (r *ID) Authorize(ctx http.Context) error {
	return nil
}

func (r *ID) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "required|string",
	}
}

func (r *ID) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
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
