package models

import (
	"github.com/goravel/framework/database/orm"
)

type Database struct {
	orm.Model
	Name     string `gorm:"unique;not null"`
	Type     string `gorm:"not null;index"`
	Host     string `gorm:"not null"`
	Port     int    `gorm:"not null"`
	Username string `gorm:"not null"`
	Password string `gorm:"default:''"`
	Remark   string `gorm:"default:''"`
}
