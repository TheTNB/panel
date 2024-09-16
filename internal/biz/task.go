package biz

import "github.com/golang-module/carbon/v2"

const (
	TaskStatusWaiting = "waiting"
	TaskStatusRunning = "running"
	TaskStatusSuccess = "finished"
	TaskStatusFailed  = "failed"
)

type Task struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"not null;index" json:"name"`
	Status    string          `gorm:"not null;default:'waiting'" json:"status"`
	Shell     string          `gorm:"not null" json:"-"`
	Log       string          `gorm:"not null" json:"log"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type TaskRepo interface {
	HasRunningTask() bool
	List(page, limit uint) ([]*Task, int64, error)
	Get(id uint) (*Task, error)
	Delete(id uint) error
}
