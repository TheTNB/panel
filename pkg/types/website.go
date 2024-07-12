package types

// WebsiteSetting 网站设置
type WebsiteSetting struct {
	Name              string   `json:"name"`
	Domains           []string `json:"domains"`
	Ports             []string `json:"ports"`
	SSLPorts          []string `json:"ssl_ports"`
	QUICPorts         []string `json:"quic_ports"`
	Root              string   `json:"root"`
	Path              string   `json:"path"`
	Index             string   `json:"index"`
	PHP               string   `json:"php"`
	OpenBasedir       bool     `json:"open_basedir"`
	SSL               bool     `json:"ssl"`
	SSLCertificate    string   `json:"ssl_certificate"`
	SSLCertificateKey string   `json:"ssl_certificate_key"`
	SSLNotBefore      string   `json:"ssl_not_before"`
	SSLNotAfter       string   `json:"ssl_not_after"`
	SSLDNSNames       []string `json:"ssl_dns_names"`
	SSLIssuer         string   `json:"ssl_issuer"`
	SSLOCSPServer     []string `json:"ssl_ocsp_server"`
	HTTPRedirect      bool     `json:"http_redirect"`
	HSTS              bool     `json:"hsts"`
	OCSP              bool     `json:"ocsp"`
	Waf               bool     `json:"waf"`
	WafMode           string   `json:"waf_mode"`
	WafCcDeny         string   `json:"waf_cc_deny"`
	WafCache          string   `json:"waf_cache"`
	Rewrite           string   `json:"rewrite"`
	Raw               string   `json:"raw"`
	Log               string   `json:"log"`
}
