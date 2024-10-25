package biz

import (
	"time"

	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/acme"
)

type Cert struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AccountID uint      `gorm:"not null" json:"account_id"` // 关联的 ACME 账户 ID
	WebsiteID uint      `gorm:"not null" json:"website_id"` // 关联的网站 ID
	DNSID     uint      `gorm:"not null" json:"dns_id"`     // 关联的 DNS ID
	Type      string    `gorm:"not null" json:"type"`       // 证书类型 (P256, P384, 2048, 3072, 4096)
	Domains   []string  `gorm:"not null;serializer:json" json:"domains"`
	AutoRenew bool      `gorm:"not null" json:"auto_renew"` // 自动续签
	CertURL   string    `gorm:"not null" json:"cert_url"`   // 证书 URL (续签时使用)
	Cert      string    `gorm:"not null" json:"cert"`       // 证书内容
	Key       string    `gorm:"not null" json:"key"`        // 私钥内容
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Website *Website     `gorm:"foreignKey:WebsiteID" json:"website"`
	Account *CertAccount `gorm:"foreignKey:AccountID" json:"account"`
	DNS     *CertDNS     `gorm:"foreignKey:DNSID" json:"dns"`
}

type CertRepo interface {
	List(page, limit uint) ([]*Cert, int64, error)
	Get(id uint) (*Cert, error)
	GetByWebsite(WebsiteID uint) (*Cert, error)
	Create(req *request.CertCreate) (*Cert, error)
	Update(req *request.CertUpdate) error
	Delete(id uint) error
	ObtainAuto(id uint) (*acme.Certificate, error)
	ObtainManual(id uint) (*acme.Certificate, error)
	Renew(id uint) (*acme.Certificate, error)
	ManualDNS(id uint) ([]acme.DNSRecord, error)
	Deploy(ID, WebsiteID uint) error
}
