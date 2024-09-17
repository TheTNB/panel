package request

import "github.com/TheTNB/panel/pkg/acme"

type CertDNSCreate struct {
	Type string        `form:"type" json:"type"`
	Name string        `form:"name" json:"name"`
	Data acme.DNSParam `form:"data" json:"data"`
}

type CertDNSUpdate struct {
	ID   uint          `form:"id" json:"id"`
	Type string        `form:"type" json:"type"`
	Name string        `form:"name" json:"name"`
	Data acme.DNSParam `form:"data" json:"data"`
}
