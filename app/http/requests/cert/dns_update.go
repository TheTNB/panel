package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"

	"github.com/TheTNB/panel/v2/pkg/acme"
)

type DNSUpdate struct {
	ID   uint          `form:"id" json:"id"`
	Type string        `form:"type" json:"type"`
	Name string        `form:"name" json:"name"`
	Data acme.DNSParam `form:"data" json:"data"`
}

func (r *DNSUpdate) Authorize(ctx http.Context) error {
	return nil
}

func (r *DNSUpdate) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id":              "required|uint|min:1|exists:cert_dns,id",
		"type":            "required|in:dnspod,aliyun,cloudflare",
		"name":            "required",
		"data":            "required",
		"data.id":         "required_if:type,dnspod",
		"data.token":      "required_if:type,dnspod",
		"data.access_key": "required_if:type,aliyun",
		"data.secret_key": "required_if:type,aliyun",
		"data.email":      "required_if:type,cloudflare",
		"data.api_key":    "required_if:type,cloudflare",
	}
}

func (r *DNSUpdate) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSUpdate) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSUpdate) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSUpdate) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
