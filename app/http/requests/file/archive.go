package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Archive struct {
	Paths []string `form:"paths" json:"paths"`
	File  string   `form:"file" json:"file"`
}

func (r *Archive) Authorize(ctx http.Context) error {
	return nil
}

func (r *Archive) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"paths":   "array",
		"paths.*": "regex:^/[a-zA-Z0-9_.@#$%-]+(\\/[a-zA-Z0-9_.@#$%-]+)*$|path_exists",
		"file":    "regex:^/[a-zA-Z0-9_.@#$%-]+(\\/[a-zA-Z0-9_.@#$%-]+)*$|path_not_exists",
	}
}

func (r *Archive) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Archive) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Archive) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
