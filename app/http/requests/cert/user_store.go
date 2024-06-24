package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UserStore struct {
	CA          string `form:"ca" json:"ca"`
	Email       string `form:"email" json:"email"`
	Kid         string `form:"kid" json:"kid"`
	HmacEncoded string `form:"hmac_encoded" json:"hmac_encoded"`
	KeyType     string `form:"key_type" json:"key_type"`
}

func (r *UserStore) Authorize(ctx http.Context) error {
	return nil
}

func (r *UserStore) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"ca":           "required|in:letsencrypt,zerossl,sslcom,google,buypass",
		"email":        "required|email",
		"kid":          "required_if:ca,sslcom,google",
		"hmac_encoded": "required_if:ca,sslcom,google",
		"key_type":     "required|in:P256,P384,2048,4096",
	}
}

func (r *UserStore) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UserStore) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UserStore) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
