package models

import "github.com/goravel/framework/support/carbon"

type User struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Username  string          `gorm:"unique;not null" json:"username"`
	Password  string          `gorm:"not null" json:"password"`
	Email     string          `gorm:"default:''" json:"email"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
