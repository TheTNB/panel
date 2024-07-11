package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type CertUpdate struct {
	ID        uint     `form:"id" json:"id"`
	Type      string   `form:"type" json:"type"`
	Domains   []string `form:"domains" json:"domains"`
	AutoRenew bool     `form:"auto_renew" json:"auto_renew"`
	UserID    uint     `form:"user_id" json:"user_id"`
	DNSID     uint     `form:"dns_id" json:"dns_id"`
	WebsiteID uint     `form:"website_id" json:"website_id"`
}

func (r *CertUpdate) Authorize(ctx http.Context) error {
	return nil
}

func (r *CertUpdate) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id":         "required|uint|min:1|exists:certs,id",
		"type":       "required|in:P256,P384,2048,4096",
		"domains":    "required|slice",
		"auto_renew": "required|bool",
		"user_id":    "required|uint|exists:cert_users,id",
		"dns_id":     "uint",
		"website_id": "uint",
	}
}

func (r *CertUpdate) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"user_id":    "uint",
		"dns_id":     "uint",
		"website_id": "uint",
	}
}

func (r *CertUpdate) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertUpdate) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertUpdate) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
