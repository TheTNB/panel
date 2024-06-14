package models

import "github.com/goravel/framework/database/orm"

type Database struct {
	orm.Model
	Name     string `gorm:"not null;unique" json:"name"`
	Type     string `gorm:"not null" json:"type"`
	Host     string `gorm:"not null" json:"host"`
	Port     int    `gorm:"not null" json:"port"`
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Remark   string `gorm:"not null" json:"remark"`
}
