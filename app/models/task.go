package models

import (
	"github.com/goravel/framework/database/orm"
)

type Task struct {
	orm.Model
	Name   string `gorm:"not null"`
	Status string `gorm:"not null;default:'waiting'"`
	Shell  string `gorm:"default:null"`
	Log    string `gorm:"default:null"`
}
