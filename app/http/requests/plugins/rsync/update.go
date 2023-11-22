package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Update struct {
	Name       string `form:"name" json:"name"`
	Path       string `form:"path" json:"path"`
	Comment    string `form:"comment" json:"comment"`
	AuthUser   string `form:"auth_user" json:"auth_user"`
	Secret     string `form:"secret" json:"secret"`
	HostsAllow string `form:"hosts_allow" json:"hosts_allow"`
}

func (r *Update) Authorize(ctx http.Context) error {
	return nil
}

func (r *Update) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "required|regex:^[a-zA-Z0-9-_]+$",
		"path":        "regex:^/[a-zA-Z0-9_-]+(\\/[a-zA-Z0-9_-]+)*$",
		"comment":     "string",
		"auth_user":   "required|regex:^[a-zA-Z0-9-_]+$",
		"secret":      "required|min_len:8",
		"hosts_allow": "string",
	}
}

func (r *Update) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Update) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Update) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
