package models

import "github.com/goravel/framework/database/orm"

type Website struct {
	orm.Model
	Name   string `gorm:"not null;unique" json:"name"`
	Status bool   `gorm:"not null;default:true" json:"status"`
	Path   string `gorm:"not null" json:"path"`
	PHP    int    `gorm:"not null" json:"php"`
	SSL    bool   `gorm:"not null" json:"ssl"`
	Remark string `gorm:"not null" json:"remark"`

	Cert *Cert `gorm:"foreignKey:WebsiteID" json:"cert"`
}
