package types

// WebsiteListen 网站监听配置
type WebsiteListen struct {
	Address string `form:"address" json:"address" validate:"required"` // 监听地址 e.g. 80 0.0.0.0:80 [::]:80
	HTTPS   bool   `form:"https" json:"https"`                         // 是否启用HTTPS
	QUIC    bool   `form:"quic" json:"quic"`                           // 是否启用QUIC
}

// WebsiteSetting 网站设置
type WebsiteSetting struct {
	ID                uint            `json:"id"`
	Name              string          `json:"name"`
	Listens           []WebsiteListen `form:"listens" json:"listens" validate:"required"`
	Domains           []string        `json:"domains"`
	Path              string          `json:"path"` // 网站目录
	Root              string          `json:"root"` // 运行目录
	Index             []string        `json:"index"`
	PHP               int             `json:"php"`
	OpenBasedir       bool            `json:"open_basedir"`
	HTTPS             bool            `json:"https"`
	SSLCertificate    string          `json:"ssl_certificate"`
	SSLCertificateKey string          `json:"ssl_certificate_key"`
	SSLNotBefore      string          `json:"ssl_not_before"`
	SSLNotAfter       string          `json:"ssl_not_after"`
	SSLDNSNames       []string        `json:"ssl_dns_names"`
	SSLIssuer         string          `json:"ssl_issuer"`
	SSLOCSPServer     []string        `json:"ssl_ocsp_server"`
	HTTPRedirect      bool            `json:"http_redirect"`
	HSTS              bool            `json:"hsts"`
	OCSP              bool            `json:"ocsp"`
	Rewrite           string          `json:"rewrite"`
	Raw               string          `json:"raw"`
	Log               string          `json:"log"`
}
