package models

import "github.com/goravel/framework/database/orm"

const (
	TaskStatusWaiting = "waiting"
	TaskStatusRunning = "running"
	TaskStatusSuccess = "finished"
	TaskStatusFailed  = "failed"
)

type Task struct {
	orm.Model
	Name   string `gorm:"not null;index" json:"name"`
	Status string `gorm:"not null;default:'waiting'" json:"status"`
	Shell  string `gorm:"not null" json:"shell"`
	Log    string `gorm:"not null" json:"log"`
}
