package biz

import (
	"time"

	"github.com/tnb-labs/panel/internal/http/request"
	"github.com/tnb-labs/panel/pkg/types"
)

type Monitor struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	Info      types.CurrentInfo `gorm:"not null;serializer:json" json:"info"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type MonitorRepo interface {
	GetSetting() (*request.MonitorSetting, error)
	UpdateSetting(setting *request.MonitorSetting) error
	Clear() error
	List(start, end time.Time) ([]*Monitor, error)
}
