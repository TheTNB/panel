package types

import "time"

type CertList struct {
	ID         uint      `json:"id"`
	AccountID  uint      `json:"account_id"`
	WebsiteID  uint      `json:"website_id"`
	DNSID      uint      `json:"dns_id"`
	Type       string    `json:"type"`
	Domains    []string  `json:"domains"`
	AutoRenew  bool      `json:"auto_renew"`
	Cert       string    `json:"cert"`
	Key        string    `json:"key"`
	NotBefore  time.Time `json:"not_before"`
	NotAfter   time.Time `json:"not_after"`
	Issuer     string    `json:"issuer"`
	OCSPServer []string  `json:"ocsp_server"`
	DNSNames   []string  `json:"dns_names"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
