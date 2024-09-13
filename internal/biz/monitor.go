package biz

import (
	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/pkg/tools"
)

type Monitor struct {
	ID        uint                 `gorm:"primaryKey" json:"id"`
	Info      tools.MonitoringInfo `gorm:"not null;serializer:json" json:"info"`
	CreatedAt carbon.DateTime      `json:"created_at"`
	UpdatedAt carbon.DateTime      `json:"updated_at"`
}
