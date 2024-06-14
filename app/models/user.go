package models

import "github.com/goravel/framework/database/orm"

type User struct {
	orm.Model
	Username string `gorm:"not null;unique" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Email    string `gorm:"not null" json:"email"`
}
