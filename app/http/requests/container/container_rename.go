package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ContainerRename struct {
	ID   string `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
}

func (r *ContainerRename) Authorize(ctx http.Context) error {
	return nil
}

func (r *ContainerRename) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id":   "required|string",
		"name": "required|string",
	}
}

func (r *ContainerRename) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ContainerRename) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ContainerRename) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ContainerRename) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
