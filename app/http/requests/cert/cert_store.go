package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/spf13/cast"
)

type CertStore struct {
	Type      string   `form:"type" json:"type"`
	Domains   []string `form:"domains" json:"domains"`
	AutoRenew bool     `form:"auto_renew" json:"auto_renew"`
	UserID    uint     `form:"user_id" json:"user_id"`
	DNSID     uint     `form:"dns_id" json:"dns_id"`
	WebsiteID uint     `form:"website_id" json:"website_id"`
}

func (r *CertStore) Authorize(ctx http.Context) error {
	return nil
}

func (r *CertStore) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"type":       "required|in:P256,P384,2048,4096",
		"domains":    "required|slice",
		"auto_renew": "required|bool",
		"user_id":    "required|uint|exists:cert_users,id",
		"dns_id":     "uint",
		"website_id": "uint",
	}
}

func (r *CertStore) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertStore) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertStore) PrepareForValidation(ctx http.Context, data validation.Data) error {
	// TODO 由于验证器 filter 标签的问题，暂时这里这样处理
	userID, exist := data.Get("user_id")
	if exist {
		err := data.Set("user_id", cast.ToUint(userID))
		if err != nil {
			return err
		}
	}
	dnsID, exist := data.Get("dns_id")
	if exist {
		err := data.Set("dns_id", cast.ToUint(dnsID))
		if err != nil {
			return err
		}

	}
	websiteID, exist := data.Get("website_id")
	if exist {
		err := data.Set("website_id", cast.ToUint(websiteID))
		if err != nil {
			return err
		}
	}

	return nil
}
