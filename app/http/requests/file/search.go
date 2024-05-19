package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Search struct {
	Path    string `form:"path" json:"path"`
	KeyWord string `form:"keyword" json:"keyword"`
}

func (r *Search) Authorize(ctx http.Context) error {
	return nil
}

func (r *Search) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"path":    "regex:^/[a-zA-Z0-9_.@#$%- []()]+(/[a-zA-Z0-9_.@#$%- []()]+)*$|path_exists",
		"keyword": "required|string",
	}
}

func (r *Search) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Search) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Search) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
