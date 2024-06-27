package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ImagePull struct {
	Name     string `form:"name" json:"name"`
	Auth     bool   `form:"auth" json:"auth"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func (r *ImagePull) Authorize(ctx http.Context) error {
	return nil
}

func (r *ImagePull) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":     "required|string",
		"auth":     "bool",
		"username": "string",
		"password": "string",
	}
}

func (r *ImagePull) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ImagePull) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ImagePull) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ImagePull) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
