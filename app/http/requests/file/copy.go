package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Copy struct {
	Old string `form:"old" json:"old"`
	New string `form:"new" json:"new"`
}

func (r *Copy) Authorize(ctx http.Context) error {
	return nil
}

func (r *Copy) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"old": "regex:^/[a-zA-Z0-9_-]+(\\/[a-zA-Z0-9_-]+)*$|path_exists",
		"new": "regex:^/[a-zA-Z0-9_-]+(\\/[a-zA-Z0-9_-]+)*$",
	}
}

func (r *Copy) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Copy) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Copy) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
