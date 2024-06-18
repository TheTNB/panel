package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Save struct {
	Path    string `form:"path" json:"path"`
	Content string `form:"content" json:"content"`
}

func (r *Save) Authorize(ctx http.Context) error {
	return nil
}

func (r *Save) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"path":    `regex:^/.*$|path_exists`,
		"content": "required|string",
	}
}

func (r *Save) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Save) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Save) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
