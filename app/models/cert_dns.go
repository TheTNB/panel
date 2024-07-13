package models

import (
	"github.com/goravel/framework/database/orm"

	"github.com/TheTNB/panel/v2/pkg/acme"
)

type CertDNS struct {
	orm.Model
	Name string        `gorm:"not null" json:"name"` // 备注名称
	Type string        `gorm:"not null" json:"type"` // DNS 提供商 (dnspod, tencent, aliyun, cloudflare)
	Data acme.DNSParam `gorm:"not null;serializer:json" json:"dns_param"`

	Certs []*Cert `gorm:"foreignKey:DNSID" json:"-"`
}
