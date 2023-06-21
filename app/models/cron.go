package models

import (
	"github.com/goravel/framework/database/orm"
)

type Cron struct {
	orm.Model
	Name   string `gorm:"not null;unique"`
	Status bool   `gorm:"not null;default:false"`
	Type   string `gorm:"not null"`
	Time   string `gorm:"not null"`
	Shell  string `gorm:"default:null"`
	Log    string `gorm:"default:null"`
}
