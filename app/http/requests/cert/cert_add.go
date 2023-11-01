package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type CertAdd struct {
	Type    string   `form:"type" json:"type"`
	Domains []string `form:"domains" json:"domains"`
	UserID  uint     `form:"user_id" json:"user_id"`
	DNSID   *uint    `form:"dns_id" json:"dns_id"`
}

func (r *CertAdd) Authorize(ctx http.Context) error {
	return nil
}

func (r *CertAdd) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"type":    "required|in:P256,P384,2048,4096",
		"domains": "required|array",
		"user_id": "required|exists:cert_users,id",
	}
}

func (r *CertAdd) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"type.required":    "类型不能为空",
		"type.in":          "类型必须为 P256, P384, 2048, 4096 中的一个",
		"domains.required": "域名不能为空",
		"domains.slice":    "域名必须为数组",
		"user_id.required": "ACME 用户 ID 不能为空",
		"user_id.exists":   "ACME 用户 ID 不存在",
	}
}

func (r *CertAdd) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertAdd) PrepareForValidation(ctx http.Context, data validation.Data) error {
	err := data.Set("user_id", uint(ctx.Request().InputInt("user_id")))
	if err != nil {
		return err
	}

	return nil
}
