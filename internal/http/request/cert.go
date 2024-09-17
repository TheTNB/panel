package request

type CertCreate struct {
	Type      string   `form:"type" json:"type"`
	Domains   []string `form:"domains" json:"domains"`
	AutoRenew bool     `form:"auto_renew" json:"auto_renew"`
	AccountID uint     `form:"account_id" json:"account_id"`
	DNSID     uint     `form:"dns_id" json:"dns_id"`
	WebsiteID uint     `form:"website_id" json:"website_id"`
}

type CertUpdate struct {
	ID        uint     `form:"id" json:"id"`
	Type      string   `form:"type" json:"type"`
	Domains   []string `form:"domains" json:"domains"`
	AutoRenew bool     `form:"auto_renew" json:"auto_renew"`
	AccountID uint     `form:"account_id" json:"account_id"`
	DNSID     uint     `form:"dns_id" json:"dns_id"`
	WebsiteID uint     `form:"website_id" json:"website_id"`
}

type CertDeploy struct {
	ID        uint `form:"id" json:"id"`
	WebsiteID uint `form:"website_id" json:"website_id"`
}
