package models

import (
	"github.com/goravel/framework/database/orm"
)

type Cert struct {
	orm.Model
	UserID    uint     `gorm:"not null" json:"user_id"`    // 关联的 ACME 用户 ID
	WebsiteID uint     `gorm:"not null" json:"website_id"` // 关联的网站 ID
	DNSID     uint     `gorm:"not null" json:"dns_id"`     // 关联的 DNS ID
	Type      string   `gorm:"not null" json:"type"`       // 证书类型 (P256, P384, 2048, 4096)
	Domains   []string `gorm:"not null;serializer:json" json:"domains"`
	AutoRenew bool     `gorm:"not null" json:"auto_renew"` // 自动续签
	CertURL   string   `gorm:"not null" json:"cert_url"`   // 证书 URL (续签时使用)
	Cert      string   `gorm:"not null" json:"cert"`       // 证书内容
	Key       string   `gorm:"not null" json:"key"`        // 私钥内容

	Website *Website  `gorm:"foreignKey:WebsiteID" json:"website"`
	User    *CertUser `gorm:"foreignKey:UserID" json:"user"`
	DNS     *CertDNS  `gorm:"foreignKey:DNSID" json:"dns"`
}
