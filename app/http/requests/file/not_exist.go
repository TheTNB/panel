package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type NotExist struct {
	Path string `form:"path" json:"path"`
}

func (r *NotExist) Authorize(ctx http.Context) error {
	return nil
}

func (r *NotExist) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"path": "regex:^/[a-zA-Z0-9_-]+(\\/[a-zA-Z0-9_-]+)*$|path_not_exists",
	}
}

func (r *NotExist) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *NotExist) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *NotExist) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
