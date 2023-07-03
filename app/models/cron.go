package models

import (
	"github.com/goravel/framework/support/carbon"
)

type Cron struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"not null;unique" json:"name"`
	Status    bool            `gorm:"not null;default:false" json:"status"`
	Type      string          `gorm:"not null" json:"type"`
	Time      string          `gorm:"not null" json:"time"`
	Shell     string          `gorm:"default:''" json:"shell"`
	Log       string          `gorm:"default:''" json:"log"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
