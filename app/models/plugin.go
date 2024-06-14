package models

import "github.com/goravel/framework/database/orm"

type Plugin struct {
	orm.Model
	Slug      string `gorm:"not null;unique" json:"slug"`
	Version   string `gorm:"not null" json:"version"`
	Show      bool   `gorm:"not null" json:"show"`
	ShowOrder int    `gorm:"not null" json:"show_order"`
}
