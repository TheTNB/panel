package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"panel/pkg/acme"
)

type DNSAdd struct {
	Type string        `form:"type" json:"type"`
	Data acme.DNSParam `form:"data" json:"data"`
}

func (r *DNSAdd) Authorize(ctx http.Context) error {
	return nil
}

func (r *DNSAdd) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"type":            "required|in:dnspod,aliyun,cloudflare",
		"data":            "required",
		"data.id":         "required_if:type,dnspod",
		"data.token":      "required_if:type,dnspod",
		"data.access_key": "required_if:type,aliyun",
		"data.secret_key": "required_if:type,aliyun",
		"data.email":      "required_if:type,cloudflare",
		"data.api_key":    "required_if:type,cloudflare",
	}
}

func (r *DNSAdd) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"type.required":            "类型不能为空",
		"type.in":                  "类型必须为 dnspod, aliyun, cloudflare 中的一个",
		"data.required":            "数据不能为空",
		"data.id.required_if":      "ID 不能为空",
		"data.token.required_if":   "Token 不能为空",
		"data.access_key.required": "Access Key 不能为空",
		"data.secret_key.required": "Secret Key 不能为空",
		"data.email.required":      "Email 不能为空",
		"data.api_key.required":    "API Key 不能为空",
	}
}

func (r *DNSAdd) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSAdd) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
