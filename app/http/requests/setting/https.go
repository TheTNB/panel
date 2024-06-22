package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Https struct {
	Https bool   `form:"https" json:"https"`
	Cert  string `form:"cert" json:"cert"`
	Key   string `form:"key" json:"key"`
}

func (r *Https) Authorize(ctx http.Context) error {
	return nil
}

func (r *Https) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"https": "bool",
		"cert":  "string",
		"key":   "string",
	}
}

func (r *Https) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Https) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Https) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
