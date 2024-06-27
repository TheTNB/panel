package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Exist struct {
	Path string `form:"path" json:"path"`
}

func (r *Exist) Authorize(ctx http.Context) error {
	return nil
}

func (r *Exist) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"path": `regex:^/.*$|path_exists`,
	}
}

func (r *Exist) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Exist) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Exist) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Exist) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
