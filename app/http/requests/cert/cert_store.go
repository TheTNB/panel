package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type CertStore struct {
	Type      string   `form:"type" json:"type"`
	Domains   []string `form:"domains" json:"domains"`
	AutoRenew bool     `form:"auto_renew" json:"auto_renew"`
	UserID    uint     `form:"user_id" json:"user_id"`
	DNSID     *uint    `form:"dns_id" json:"dns_id"`
	WebsiteID *uint    `form:"website_id" json:"website_id"`
}

func (r *CertStore) Authorize(ctx http.Context) error {
	return nil
}

func (r *CertStore) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"type":       "required|in:P256,P384,2048,4096",
		"domains":    "required|array",
		"auto_renew": "required|bool",
		"user_id":    "required|exists:cert_users,id",
	}
}

func (r *CertStore) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertStore) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertStore) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
