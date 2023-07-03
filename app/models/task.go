package models

import (
	"github.com/goravel/framework/support/carbon"
)

const (
	TaskStatusWaiting = "waiting"
	TaskStatusRunning = "running"
	TaskStatusSuccess = "finished"
	TaskStatusFailed  = "failed"
)

type Task struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"not null" json:"name"`
	Status    string          `gorm:"not null;default:'waiting'" json:"status"`
	Shell     string          `gorm:"default:''" json:"shell"`
	Log       string          `gorm:"default:''" json:"log"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
