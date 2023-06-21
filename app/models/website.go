package models

import (
	"github.com/goravel/framework/database/orm"
)

type Website struct {
	orm.Model
	Name   string `gorm:"unique;not null"`
	Status bool   `gorm:"default:true;not null;index"`
	Path   string `gorm:"not null"`
	Php    int    `gorm:"default:0;not null;index"`
	Ssl    bool   `gorm:"default:false;not null;index"`
	Remark string `gorm:"default:null"`
}
