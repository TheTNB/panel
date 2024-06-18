package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Move struct {
	Source string `form:"source" json:"source"`
	Target string `form:"target" json:"target"`
}

func (r *Move) Authorize(ctx http.Context) error {
	return nil
}

func (r *Move) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"source": `regex:^/.*$|path_exists`,
		"target": `regex:^/.*$`,
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
