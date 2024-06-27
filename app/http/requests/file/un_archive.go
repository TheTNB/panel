package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UnArchive struct {
	File string `form:"file" json:"file"`
	Path string `form:"path" json:"path"`
}

func (r *UnArchive) Authorize(ctx http.Context) error {
	return nil
}

func (r *UnArchive) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"file": `regex:^/.*$|path_exists`,
		"path": `regex:^/.*$`,
	}
}

func (r *UnArchive) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UnArchive) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UnArchive) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UnArchive) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
