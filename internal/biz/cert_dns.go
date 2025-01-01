package biz

import (
	"time"

	"github.com/tnb-labs/panel/internal/http/request"
	"github.com/tnb-labs/panel/pkg/acme"
)

type CertDNS struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Name      string        `gorm:"not null" json:"name"` // 备注名称
	Type      string        `gorm:"not null" json:"type"` // DNS 提供商 (tencent, aliyun, cloudflare)
	Data      acme.DNSParam `gorm:"not null;serializer:json" json:"dns_param"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`

	Certs []*Cert `gorm:"foreignKey:DNSID" json:"-"`
}

type CertDNSRepo interface {
	List(page, limit uint) ([]*CertDNS, int64, error)
	Get(id uint) (*CertDNS, error)
	Create(req *request.CertDNSCreate) (*CertDNS, error)
	Update(req *request.CertDNSUpdate) error
	Delete(id uint) error
}
