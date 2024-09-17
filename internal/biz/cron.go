package biz

import (
	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/internal/http/request"
)

type Cron struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"not null;unique" json:"name"`
	Status    bool            `gorm:"not null" json:"status"`
	Type      string          `gorm:"not null" json:"type"`
	Time      string          `gorm:"not null" json:"time"`
	Shell     string          `gorm:"not null" json:"shell"`
	Log       string          `gorm:"not null" json:"log"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type CronRepo interface {
	Count() (int64, error)
	List(page, limit uint) ([]*Cron, int64, error)
	Get(id uint) (*Cron, error)
	Create(req *request.CronCreate) error
	Update(req *request.CronUpdate) error
	Delete(id uint) error
	Status(id uint, status bool) error
	Log(id uint) (string, error)
}
