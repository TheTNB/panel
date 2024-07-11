package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type RestoreBackup struct {
	ID   uint   `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
}

func (r *RestoreBackup) Authorize(ctx http.Context) error {
	return nil
}

func (r *RestoreBackup) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id":   "required|exists:websites,id",
		"name": "required|string",
	}
}

func (r *RestoreBackup) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "uint",
	}
}

func (r *RestoreBackup) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *RestoreBackup) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *RestoreBackup) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
