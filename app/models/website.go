package models

import (
	"github.com/goravel/framework/support/carbon"
)

type Website struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `json:"name"`
	Status    bool            `gorm:"default:true" json:"status"`
	Path      string          `json:"path"`
	Php       int             `gorm:"default:0;not null;index" json:"php"`
	Ssl       bool            `gorm:"default:false;not null;index" json:"ssl"`
	Remark    string          `gorm:"default:''" json:"remark"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`

	Cert *Cert `gorm:"foreignKey:WebsiteID" json:"cert"`
}
