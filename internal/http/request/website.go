package request

type WebsiteDefaultConfig struct {
	Index string `json:"index" form:"index" validate:"required"`
	Stop  string `json:"stop" form:"stop" validate:"required"`
}

type WebsiteCreate struct {
	Name       string   `form:"name" json:"name" validate:"required"`
	Domains    []string `form:"domains" json:"domains" validate:"required"`
	Ports      []uint   `form:"ports" json:"ports" validate:"required"`
	Path       string   `form:"path" json:"path"`
	PHP        int      `form:"php" json:"php"`
	DB         bool     `form:"db" json:"db"`
	DBType     string   `form:"db_type" json:"db_type"`
	DBName     string   `form:"db_name" json:"db_name"`
	DBUser     string   `form:"db_user" json:"db_user"`
	DBPassword string   `form:"db_password" json:"db_password"`
}

type WebsiteDelete struct {
	ID   uint `form:"id" json:"id" validate:"required"`
	Path bool `form:"path" json:"path"`
	DB   bool `form:"db" json:"db"`
}

type WebsiteUpdate struct {
	ID                uint     `form:"id" json:"id" validate:"required"`
	Domains           []string `form:"domains" json:"domains" validate:"required"`
	Ports             []uint   `form:"ports" json:"ports" validate:"required"`
	SSLPorts          []uint   `form:"ssl_ports" json:"ssl_ports" validate:"required"`
	QUICPorts         []uint   `form:"quic_ports" json:"quic_ports" validate:"required"`
	OCSP              bool     `form:"ocsp" json:"ocsp"`
	HSTS              bool     `form:"hsts" json:"hsts"`
	SSL               bool     `form:"ssl" json:"ssl"`
	HTTPRedirect      bool     `form:"http_redirect" json:"http_redirect"`
	OpenBasedir       bool     `form:"open_basedir" json:"open_basedir"`
	Index             string   `form:"index" json:"index" validate:"required"`
	Path              string   `form:"path" json:"path" validate:"required"`
	Root              string   `form:"root" json:"root" validate:"required"`
	Raw               string   `form:"raw" json:"raw"`
	Rewrite           string   `form:"rewrite" json:"rewrite"`
	PHP               int      `form:"php" json:"php"`
	SSLCertificate    string   `form:"ssl_certificate" json:"ssl_certificate"`
	SSLCertificateKey string   `form:"ssl_certificate_key" json:"ssl_certificate_key"`
}

type WebsiteUpdateRemark struct {
	ID     uint   `form:"id" json:"id" validate:"required"`
	Remark string `form:"remark" json:"remark"`
}

type WebsiteUpdateStatus struct {
	ID     uint `json:"id" form:"id" validate:"required"`
	Status bool `json:"status" form:"status"`
}
