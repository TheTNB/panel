package request

type CertUpload struct {
	Cert string `form:"cert" json:"cert" validate:"required"`
	Key  string `form:"key" json:"key" validate:"required"`
}

type CertCreate struct {
	Type      string   `form:"type" json:"type" validate:"required|in:P256,P384,2048,3072,4096"`
	Domains   []string `form:"domains" json:"domains" validate:"required|isSlice"`
	AutoRenew bool     `form:"auto_renew" json:"auto_renew"`
	AccountID uint     `form:"account_id" json:"account_id"`
	DNSID     uint     `form:"dns_id" json:"dns_id"`
	WebsiteID uint     `form:"website_id" json:"website_id"`
}

type CertUpdate struct {
	ID        uint     `form:"id" json:"id" validate:"required|exists:certs,id"`
	Type      string   `form:"type" json:"type" validate:"required|in:P256,P384,2048,3072,4096"`
	Domains   []string `form:"domains" json:"domains" validate:"required|isSlice"`
	Cert      string   `form:"cert" json:"cert"`
	Key       string   `form:"key" json:"key"`
	AutoRenew bool     `form:"auto_renew" json:"auto_renew"`
	AccountID uint     `form:"account_id" json:"account_id"`
	DNSID     uint     `form:"dns_id" json:"dns_id"`
	WebsiteID uint     `form:"website_id" json:"website_id"`
}

type CertDeploy struct {
	ID        uint `form:"id" json:"id" validate:"required|exists:certs,id"`
	WebsiteID uint `form:"website_id" json:"website_id" validate:"required|exists:websites,id"`
}
