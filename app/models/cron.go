package models

import "github.com/goravel/framework/database/orm"

type Cron struct {
	orm.Model
	Name   string `gorm:"not null;unique" json:"name"`
	Status bool   `gorm:"not null" json:"status"`
	Type   string `gorm:"not null" json:"type"`
	Time   string `gorm:"not null" json:"time"`
	Shell  string `gorm:"not null" json:"shell"`
	Log    string `gorm:"not null" json:"log"`
}
