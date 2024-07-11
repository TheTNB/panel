package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Delete struct {
	ID   uint `form:"id" json:"id"`
	Path bool `form:"path" json:"path"`
	DB   bool `form:"db" json:"db"`
}

func (r *Delete) Authorize(ctx http.Context) error {
	return nil
}

func (r *Delete) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id":   "required|exists:websites,id",
		"path": "bool",
		"db":   "bool",
	}
}

func (r *Delete) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "uint",
	}
}

func (r *Delete) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Delete) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Delete) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
