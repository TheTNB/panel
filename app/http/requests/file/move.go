package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Move struct {
	Old string `form:"old" json:"old"`
	New string `form:"new" json:"new"`
}

func (r *Move) Authorize(ctx http.Context) error {
	return nil
}

func (r *Move) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"old": "regex:^/[a-zA-Z0-9_-]+(\\/[a-zA-Z0-9_-]+)*$|path_exists",
		"new": "regex:^/[a-zA-Z0-9_-]+(\\/[a-zA-Z0-9_-]+)*$|path_not_exists",
	}
}

func (r *Move) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Move) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Move) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
