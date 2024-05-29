package internal

import (
	requests "github.com/TheTNB/panel/app/http/requests/website"
	"github.com/TheTNB/panel/app/models"
)

type Website interface {
	List(page int, limit int) (int64, []models.Website, error)
	Add(website PanelWebsite) (models.Website, error)
	SaveConfig(config requests.SaveConfig) error
	Delete(id uint) error
	GetConfig(id uint) (WebsiteSetting, error)
	GetConfigByName(name string) (WebsiteSetting, error)
	GetIDByName(name string) (uint, error)
}

type PanelWebsite struct {
	Name       string   `json:"name"`
	Status     bool     `json:"status"`
	Domains    []string `json:"domains"`
	Ports      []uint   `json:"ports"`
	Path       string   `json:"path"`
	Php        string   `json:"php"`
	Ssl        bool     `json:"ssl"`
	Remark     string   `json:"remark"`
	Db         bool     `json:"db"`
	DbType     string   `json:"db_type"`
	DbName     string   `json:"db_name"`
	DbUser     string   `json:"db_user"`
	DbPassword string   `json:"db_password"`
}

// WebsiteSetting 网站设置
type WebsiteSetting struct {
	Name              string   `json:"name"`
	Domains           []string `json:"domains"`
	Ports             []uint   `json:"ports"`
	Root              string   `json:"root"`
	Path              string   `json:"path"`
	Index             string   `json:"index"`
	Php               string   `json:"php"`
	OpenBasedir       bool     `json:"open_basedir"`
	Ssl               bool     `json:"ssl"`
	SslCertificate    string   `json:"ssl_certificate"`
	SslCertificateKey string   `json:"ssl_certificate_key"`
	SslNotBefore      string   `json:"ssl_not_before"`
	SslNotAfter       string   `json:"ssl_not_after"`
	SSlDNSNames       []string `json:"ssl_dns_names"`
	SslIssuer         string   `json:"ssl_issuer"`
	SslOCSPServer     []string `json:"ssl_ocsp_server"`
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
