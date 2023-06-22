package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`

	CaptchaID string `json:"captcha_id" form:"captcha_id"`
	Captcha   string `json:"captcha" form:"captcha"`
}

func (r *LoginRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *LoginRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"username": "required",
		"password": "required|min_len:8",
		"captcha":  "captcha:true",
	}
}

func (r *LoginRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"username.required": "登录名不能为空",
		"password.required": "密码不能为空",
		"password.min_len":  "密码长度不能小于 8 位",
		"captcha.captcha":   "验证码错误",
	}
}

func (r *LoginRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *LoginRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
