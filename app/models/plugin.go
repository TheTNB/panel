package models

import (
	"github.com/goravel/framework/database/orm"
)

type Plugin struct {
	orm.Model
	Slug      string `gorm:"unique;not null"`
	Version   string `gorm:"not null"`
	Show      bool   `gorm:"default:false;not null"`
	ShowOrder int    `gorm:"default:0;not null"`
}
