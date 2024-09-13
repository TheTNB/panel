package biz

import (
	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/pkg/acme"
)

type CertDNS struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"not null" json:"name"` // 备注名称
	Type      string          `gorm:"not null" json:"type"` // DNS 提供商 (dnspod, tencent, aliyun, cloudflare)
	Data      acme.DNSParam   `gorm:"not null;serializer:json" json:"dns_param"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`

	Certs []*Cert `gorm:"foreignKey:DNSID" json:"-"`
}
