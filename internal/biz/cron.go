package biz

import "github.com/golang-module/carbon/v2"

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
}
