package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Add struct {
	Name       string   `form:"name" json:"name"`
	Domains    []string `form:"domains" json:"domains"`
	Ports      []uint   `form:"ports" json:"ports"`
	Path       string   `form:"path" json:"path"`
	PHP        string   `form:"php" json:"php"`
	DB         bool     `form:"db" json:"db"`
	DBType     string   `form:"db_type" json:"db_type"`
	DBName     string   `form:"db_name" json:"db_name"`
	DBUser     string   `form:"db_user" json:"db_user"`
	DBPassword string   `form:"db_password" json:"db_password"`
}

func (r *Add) Authorize(ctx http.Context) error {
	return nil
}

func (r *Add) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "required|regex:^[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)*$|not_exists:websites,name|not_in:phpmyadmin,mysql,panel,ssh",
		"domains":     "required|slice",
		"ports":       "required|slice",
		"path":        `regex:^/.*$`,
		"php":         "required",
		"db":          "bool",
		"db_type":     "required_if:db,true|in:0,mysql,postgresql",
		"db_name":     "required_if:db,true|regex:^[a-zA-Z0-9_-]+$",
		"db_user":     "required_if:db,true|regex:^[a-zA-Z0-9_-]+$",
		"db_password": "required_if:db,true|min_len:8",
	}
}

func (r *Add) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Add) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Add) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Add) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
