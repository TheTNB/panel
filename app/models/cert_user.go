package models

import "github.com/goravel/framework/database/orm"

type CertUser struct {
	orm.Model
	Email       string `gorm:"not null" json:"email"`
	CA          string `gorm:"not null" json:"ca"` // CA 提供商 (letsencrypt, zerossl, sslcom, google, buypass)
	Kid         string `gorm:"not null" json:"kid"`
	HmacEncoded string `gorm:"not null" json:"hmac_encoded"`
	PrivateKey  string `gorm:"not null" json:"private_key"`
	KeyType     string `gorm:"not null" json:"key_type"`

	Certs []*Cert `gorm:"foreignKey:UserID" json:"-"`
}
