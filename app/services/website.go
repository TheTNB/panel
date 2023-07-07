// Package services 网站服务
package services

import "panel/app/models"

type Website interface {
	List() ([]models.Website, error)
}

type PanelWebsite struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	Path       string `json:"path"`
	Php        string `json:"php"`
	Ssl        string `json:"ssl"`
	Remark     string `json:"remark"`
	Db         bool   `json:"db"`
	DbType     string `json:"db_type"`
	DbName     string `json:"db_name"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
}

// WebsiteSetting 网站设置
type WebsiteSetting struct {
	Name              string   `json:"name"`
	Ports             []string `json:"ports"`
	Domains           []string `json:"domains"`
	Root              string   `json:"root"`
	Path              string   `json:"path"`
	Index             string   `json:"index"`
	OpenBasedir       bool     `json:"open_basedir"`
	Ssl               bool     `json:"ssl"`
	SslCertificate    string   `json:"ssl_certificate"`
	SslCertificateKey string   `json:"ssl_certificate_key"`
	HttpRedirect      bool     `json:"http_redirect"`
	Hsts              bool     `json:"hsts"`
	Waf               bool     `json:"waf"`
	WafMode           string   `json:"waf_mode"`
	WafCcDeny         string   `json:"waf_cc_deny"`
	WafCache          string   `json:"waf_cache"`
	Rewrite           string   `json:"rewrite"`
	Raw               string   `json:"raw"`
	Log               string   `json:"log"`
}
