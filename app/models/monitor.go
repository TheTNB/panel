package models

import (
	"panel/packages/helpers"

	"github.com/goravel/framework/database/orm"
)

type Monitor struct {
	orm.Model
	Info helpers.MonitoringInfo `gorm:"type:json;serializer:json"`
}
