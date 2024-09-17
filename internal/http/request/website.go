package request

import "net/http"

type WebsiteDefaultConfig struct {
	Index string `json:"index" form:"index"`
	Stop  string `json:"stop" form:"stop"`
}

type WebsiteCreate struct {
	Name       string   `form:"name" json:"name"`
	Domains    []string `form:"domains" json:"domains"`
	Ports      []uint   `form:"ports" json:"ports"`
	Path       string   `form:"path" json:"path"`
	PHP        string   `form:"php" json:"php"`
	DB         bool     `form:"db" json:"db"`
	DBType     string   `form:"db_type" json:"db_type"`
	DBName     string   `form:"db_name" json:"db_name"`
	DBUser     string   `form:"db_user" json:"db_user"`
	DBPassword string   `form:"db_password" json:"db_password"`
}

type WebsiteDelete struct {
	ID   uint `form:"id" json:"id"`
	Path bool `form:"path" json:"path"`
	DB   bool `form:"db" json:"db"`
}

type WebsiteUpdate struct {
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

func (r *WebsiteUpdate) Prepare(_ *http.Request) error {
	if r.WafMode == "" {
		r.WafMode = "DYNAMIC"
	}
	if r.WafCcDeny == "" {
		r.WafCcDeny = "rate=1000r/m duration=60m"
	}
	if r.WafCache == "" {
		r.WafCache = "capacity=50"
	}

	return nil
}

type WebsiteUpdateRemark struct {
	ID     uint   `form:"id" json:"id"`
	Remark string `form:"remark" json:"remark"`
}

type WebsiteUpdateStatus struct {
	ID     uint `json:"id" form:"id"`
	Status bool `json:"status" form:"status"`
}
