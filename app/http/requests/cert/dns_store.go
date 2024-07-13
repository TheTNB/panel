package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"

	"github.com/TheTNB/panel/v2/pkg/acme"
)

type DNSStore struct {
	Type string        `form:"type" json:"type"`
	Name string        `form:"name" json:"name"`
	Data acme.DNSParam `form:"data" json:"data"`
}

func (r *DNSStore) Authorize(ctx http.Context) error {
	return nil
}

func (r *DNSStore) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"type":            "required|in:dnspod,tencent,aliyun,cloudflare",
		"name":            "required",
		"data":            "required",
		"data.id":         "required_if:type,dnspod",
		"data.token":      "required_if:type,dnspod",
		"data.access_key": "required_if:type,aliyun,tencent",
		"data.secret_key": "required_if:type,aliyun,tencent",
		"data.api_key":    "required_if:type,cloudflare",
	}
}

func (r *DNSStore) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSStore) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSStore) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DNSStore) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
