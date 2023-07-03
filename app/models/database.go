package models

import "github.com/goravel/framework/support/carbon"

type Database struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"unique;not null" json:"name"`
	Type      string          `gorm:"not null;index" json:"type"`
	Host      string          `gorm:"not null" json:"host"`
	Port      int             `gorm:"not null" json:"port"`
	Username  string          `gorm:"not null" json:"username"`
	Password  string          `gorm:"default:''" json:"password"`
	Remark    string          `gorm:"default:''" json:"remark"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
