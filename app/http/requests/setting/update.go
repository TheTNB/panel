package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Update struct {
	Name        string `form:"name" json:"name"`
	Language    string `form:"language" json:"language"`
	Port        uint   `form:"port" json:"port"`
	BackupPath  string `form:"backup_path" json:"backup_path"`
	WebsitePath string `form:"website_path" json:"website_path"`
	Entrance    string `form:"entrance" json:"entrance"`
	UserName    string `form:"username" json:"username"`
	Email       string `form:"email" json:"email"`
	Password    string `form:"password" json:"password"`
}

func (r *Update) Authorize(ctx http.Context) error {
	return nil
}

func (r *Update) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":         "required|string:2,20",
		"language":     "required|in:zh_CN,en",
		"port":         "required|int:1000,65535",
		"backup_path":  "required|string:2,255",
		"website_path": "required|string:2,255",
		"entrance":     `required|regex:^/(\w+)?$|not_in:/api`,
		"username":     "required|string:2,20",
		"email":        "required|email",
		"password":     "string:8,255",
	}
}

func (r *Update) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"port": "uint",
	}
}

func (r *Update) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"port.int":        "port 值必须是一个整数且在 1000 - 65535 之间",
		"password.string": "password 必须是一个字符串且长度在 8 - 255 之间",
	}
}

func (r *Update) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Update) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
