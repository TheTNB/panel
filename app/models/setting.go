package models

import "github.com/goravel/framework/support/carbon"

type Setting struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Key       string          `gorm:"unique;not null" json:"key"`
	Value     string          `gorm:"default:''" json:"value"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
