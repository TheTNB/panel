package biz

import "github.com/golang-module/carbon/v2"

type CertUser struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Email       string          `gorm:"not null" json:"email"`
	CA          string          `gorm:"not null" json:"ca"` // CA 提供商 (letsencrypt, zerossl, sslcom, google, buypass)
	Kid         string          `gorm:"not null" json:"kid"`
	HmacEncoded string          `gorm:"not null" json:"hmac_encoded"`
	PrivateKey  string          `gorm:"not null" json:"private_key"`
	KeyType     string          `gorm:"not null" json:"key_type"`
	CreatedAt   carbon.DateTime `json:"created_at"`
	UpdatedAt   carbon.DateTime `json:"updated_at"`

	Certs []*Cert `gorm:"foreignKey:UserID" json:"-"`
}
