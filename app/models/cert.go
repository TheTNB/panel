package models

import (
	"github.com/goravel/framework/support/carbon"
)

type Cert struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"default:null" json:"user_id"`              // 关联的 ACME 用户 ID
	WebsiteID *uint           `gorm:"default:null" json:"website_id"`           // 关联的网站 ID
	DNSID     *uint           `gorm:"column:dns_id;default:null" json:"dns_id"` // 关联的 DNS ID
	CronID    *uint           `gorm:"default:null" json:"cron_id"`              // 关联的计划任务 ID
	Type      string          `gorm:"not null" json:"type"`                     // 证书类型 (P256, P384, 2048, 4096)
	Domains   []string        `gorm:"type:json;serializer:json" json:"domains"`
	CertURL   *string         `gorm:"default:null" json:"cert_url"` // 证书 URL (续签时使用)
	Cert      string          `gorm:"default:null" json:"cert"`     // 证书内容
	Key       string          `gorm:"default:null" json:"key"`      // 私钥内容
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`

	Website *Website  `gorm:"foreignKey:WebsiteID" json:"website"`
	User    *CertUser `gorm:"foreignKey:UserID" json:"user"`
	DNS     *CertDNS  `gorm:"foreignKey:DNSID" json:"dns"`
	Cron    *Cron     `gorm:"foreignKey:CronID" json:"cron"`
}
