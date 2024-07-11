package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type CertDeploy struct {
	ID        uint `form:"id" json:"id"`
	WebsiteID uint `form:"website_id" json:"website_id"`
}

func (r *CertDeploy) Authorize(ctx http.Context) error {
	return nil
}

func (r *CertDeploy) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id":         "required|uint|min:1|exists:certs,id",
		"website_id": "required|uint|min:1|exists:websites,id",
	}
}

func (r *CertDeploy) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"id":         "uint",
		"website_id": "uint",
	}
}

func (r *CertDeploy) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertDeploy) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CertDeploy) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
