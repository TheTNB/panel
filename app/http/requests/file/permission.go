package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Permission struct {
	Path  string `form:"path" json:"path"`
	Mode  uint   `form:"mode" json:"mode" filter:"uint"`
	Owner string `form:"owner" json:"owner"`
	Group string `form:"group" json:"group"`
}

func (r *Permission) Authorize(ctx http.Context) error {
	return nil
}

func (r *Permission) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"path":  `regex:^/.*$|path_exists`,
		"mode":  "regex:^[0-7]{3}$|uint",
		"owner": "regex:^[a-zA-Z0-9_-]+$",
		"group": "regex:^[a-zA-Z0-9_-]+$",
	}
}

func (r *Permission) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Permission) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Permission) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
