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
		"source": `regex:^/[a-zA-Z0-9_.@#$%\-\s\[\]()]+(/[a-zA-Z0-9_.@#$%\-\s\[\]()]+)*$|path_exists`,
		"target": `regex:^/[a-zA-Z0-9_.@#$%\-\s\[\]()]+(/[a-zA-Z0-9_.@#$%\-\s\[\]()]+)*$`,
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
