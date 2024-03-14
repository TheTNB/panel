package requests

import (
	"mime/multipart"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Upload struct {
	Path string                `form:"path" json:"path"`
	File *multipart.FileHeader `form:"file" json:"file"`
}

func (r *Upload) Authorize(ctx http.Context) error {
	return nil
}

func (r *Upload) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"path": "regex:^/[a-zA-Z0-9_.@#$%-]+(\\/[a-zA-Z0-9_.@#$%-]+)*$",
		"file": "required",
	}
}

func (r *Upload) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Upload) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Upload) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
