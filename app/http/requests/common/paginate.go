package commonrequests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Paginate struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}

func (r *Paginate) Authorize(ctx http.Context) error {
	return nil
}

func (r *Paginate) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"page":  "required|int|min:1",
		"limit": "required|int|min:1",
	}
}

func (r *Paginate) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"page":  "int",
		"limit": "int",
	}
}

func (r *Paginate) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Paginate) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Paginate) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
