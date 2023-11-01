package models

import (
	"github.com/goravel/framework/support/carbon"
)

type CertUser struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Email       string          `gorm:"not null" json:"email"`
	CA          string          `gorm:"not null" json:"ca"` // CA 提供商 (letsencrypt, zerossl, sslcom, google, buypass)
	Kid         *string         `gorm:"default:null" json:"kid"`
	HmacEncoded *string         `gorm:"default:null" json:"hmac_encoded"`
	PrivateKey  string          `gorm:"not null" json:"private_key"`
	KeyType     string          `gorm:"not null" json:"key_type"`
	CreatedAt   carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt   carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`

	Certs []*Cert `gorm:"foreignKey:UserID" json:"certs"`
}
