package biz

import "github.com/golang-module/carbon/v2"

type Database struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"not null;unique" json:"name"`
	Type      string          `gorm:"not null" json:"type"`
	Host      string          `gorm:"not null" json:"host"`
	Port      int             `gorm:"not null" json:"port"`
	Username  string          `gorm:"not null" json:"username"`
	Password  string          `gorm:"not null" json:"password"`
	Remark    string          `gorm:"not null" json:"remark"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}
