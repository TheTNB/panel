package models

import "github.com/goravel/framework/support/carbon"

type Plugin struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Slug      string          `gorm:"unique;not null" json:"slug"`
	Version   string          `gorm:"not null" json:"version"`
	Show      bool            `gorm:"default:false;not null" json:"show"`
	ShowOrder int             `gorm:"default:0;not null" json:"show_order"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
