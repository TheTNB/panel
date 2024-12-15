package request

import "github.com/TheTNB/panel/pkg/acme"

type CertDNSCreate struct {
	Type string        `form:"type" json:"type" validate:"required"`
	Name string        `form:"name" json:"name" validate:"required"`
	Data acme.DNSParam `form:"data" json:"data"`
}

type CertDNSUpdate struct {
	ID   uint          `form:"id" json:"id" validate:"required|exists:cert_dns,id"`
	Type string        `form:"type" json:"type" validate:"required"`
	Name string        `form:"name" json:"name" validate:"required"`
	Data acme.DNSParam `form:"data" json:"data"`
}
