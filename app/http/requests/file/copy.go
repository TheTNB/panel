package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Copy struct {
	Source string `form:"source" json:"source"`
	Target string `form:"target" json:"target"`
}

func (r *Copy) Authorize(ctx http.Context) error {
	return nil
}

func (r *Copy) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"source": `regex:^/.*$|path_exists`,
		"target": `regex:^/.*$`,
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
