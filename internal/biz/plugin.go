package biz

import (
	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/pkg/types"
)

type Plugin struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Slug      string          `gorm:"not null;unique" json:"slug"`
	Version   string          `gorm:"not null" json:"version"`
	Show      bool            `gorm:"not null" json:"show"`
	ShowOrder int             `gorm:"not null" json:"show_order"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type PluginRepo interface {
	All() []*types.Plugin
	Installed() ([]*Plugin, error)
	Get(slug string) (*types.Plugin, error)
	GetInstalled(slug string) (*Plugin, error)
	GetInstalledAll(cond ...string) ([]*Plugin, error)
	IsInstalled(cond ...string) (bool, error)
	Install(slug string) error
	Uninstall(slug string) error
	Update(slug string) error
}
