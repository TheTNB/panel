package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"default:null"`
	orm.SoftDeletes
}
