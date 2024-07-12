package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type SaveConfig struct {
	ID                uint     `form:"id" json:"id"`
	Domains           []string `form:"domains" json:"domains"`
	Ports             []uint   `form:"ports" json:"ports"`
	SSLPorts          []uint   `form:"ssl_ports" json:"ssl_ports"`
	QUICPorts         []uint   `form:"quic_ports" json:"quic_ports"`
	OCSP              bool     `form:"ocsp" json:"ocsp"`
	HSTS              bool     `form:"hsts" json:"hsts"`
	SSL               bool     `form:"ssl" json:"ssl"`
	HTTPRedirect      bool     `form:"http_redirect" json:"http_redirect"`
	OpenBasedir       bool     `form:"open_basedir" json:"open_basedir"`
	Waf               bool     `form:"waf" json:"waf"`
	WafCache          string   `form:"waf_cache" json:"waf_cache"`
	WafMode           string   `form:"waf_mode" json:"waf_mode"`
	WafCcDeny         string   `form:"waf_cc_deny" json:"waf_cc_deny"`
	Index             string   `form:"index" json:"index"`
	Path              string   `form:"path" json:"path"`
	Root              string   `form:"root" json:"root"`
	Raw               string   `form:"raw" json:"raw"`
	Rewrite           string   `form:"rewrite" json:"rewrite"`
	PHP               int      `form:"php" json:"php"`
	SSLCertificate    string   `form:"ssl_certificate" json:"ssl_certificate"`
	SSLCertificateKey string   `form:"ssl_certificate_key" json:"ssl_certificate_key"`
}

func (r *SaveConfig) Authorize(ctx http.Context) error {
	return nil
}

func (r *SaveConfig) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id":                  "required|exists:websites,id",
		"domains":             "required|slice",
		"ports":               "required|slice",
		"ssl_ports":           "slice|not_in:80",
		"quic_ports":          "slice|not_in:80",
		"ocsp":                "bool",
		"hsts":                "bool",
		"ssl":                 "bool",
		"http_redirect":       "bool",
		"open_basedir":        "bool",
		"waf":                 "bool",
		"waf_cache":           "required|string",
		"waf_mode":            "required|string",
		"waf_cc_deny":         "required|string",
		"index":               "required|string",
		"path":                "required|string",
		"root":                "required|string",
		"raw":                 "required|string",
		"rewrite":             "string",
		"php":                 "int",
		"ssl_certificate":     "required_if:ssl,true",
		"ssl_certificate_key": "required_if:ssl,true",
	}
}

func (r *SaveConfig) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"id":  "uint",
		"php": "int",
	}
}

func (r *SaveConfig) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *SaveConfig) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *SaveConfig) PrepareForValidation(ctx http.Context, data validation.Data) error {
	_, exist := data.Get("waf_mode")
	if !exist {
		if err := data.Set("waf_mode", "DYNAMIC"); err != nil {
			return err
		}
	}
	_, exist = data.Get("waf_cc_deny")
	if !exist {
		if err := data.Set("waf_cc_deny", "rate=1000r/m duration=60m"); err != nil {
			return err
		}
	}
	_, exist = data.Get("waf_cache")
	if !exist {
		if err := data.Set("waf_cache", "capacity=50"); err != nil {
			return err
		}
	}

	return nil
}
