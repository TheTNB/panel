package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type DeleteBackup struct {
	Name string `form:"name" json:"name"`
}

func (r *DeleteBackup) Authorize(ctx http.Context) error {
	return nil
}

func (r *DeleteBackup) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name": "required|string",
	}
}

func (r *DeleteBackup) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DeleteBackup) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DeleteBackup) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
