package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UserAdd struct {
	CA          string `form:"ca" json:"ca"`
	Email       string `form:"email" json:"email"`
	Kid         string `form:"kid" json:"kid"`
	HmacEncoded string `form:"hmac_encoded" json:"hmac_encoded"`
	KeyType     string `form:"key_type" json:"key_type"`
}

func (r *UserAdd) Authorize(ctx http.Context) error {
	return nil
}

func (r *UserAdd) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"ca":           "required|in:letsencrypt,zerossl,sslcom,google,buypass",
		"email":        "required|email",
		"kid":          "required_unless:ca,letsencrypt,buypass",
		"hmac_encoded": "required_unless:ca,letsencrypt,buypass",
		"key_type":     "required|in:P256,P384,2048,4096",
	}
}

func (r *UserAdd) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"ca.required":           "CA 不能为空",
		"ca.in":                 "CA 必须为 letsencrypt, zerossl, sslcom, google, buypass 中的一个",
		"email.required":        "邮箱不能为空",
		"email.email":           "邮箱格式不正确",
		"kid.required_unless":   "KID 不能为空",
		"hmac_encoded.required": "HMAC Encoded 不能为空",
		"key_type.required":     "密钥类型不能为空",
		"key_type.in":           "密钥类型必须为 P256, P384, 2048, 4096 中的一个",
	}
}

func (r *UserAdd) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UserAdd) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
