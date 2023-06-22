package models

import (
	"github.com/goravel/framework/database/orm"
)

const (
	TaskStatusWaiting = "waiting"
	TaskStatusRunning = "running"
	TaskStatusSuccess = "finished"
	TaskStatusFailed  = "failed"
)

type Task struct {
	orm.Model
	Name   string `gorm:"not null"`
	Status string `gorm:"not null;default:'waiting'"`
	Shell  string `gorm:"default:''"`
	Log    string `gorm:"default:''"`
}
