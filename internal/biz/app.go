package biz

import (
	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/pkg/api"
)

type App struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Slug      string          `gorm:"not null;unique" json:"slug"`
	Version   string          `gorm:"not null" json:"version"`
	Show      bool            `gorm:"not null" json:"show"`
	ShowOrder int             `gorm:"not null" json:"show_order"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type AppRepo interface {
	All() api.Apps
	Get(slug string) (*api.App, error)
	Installed() ([]*App, error)
	GetInstalled(slug string) (*App, error)
	GetInstalledAll(cond ...string) ([]*App, error)
	IsInstalled(cond ...string) (bool, error)
	Install(slug string) error
	Uninstall(slug string) error
	Update(slug string) error
	UpdateShow(slug string, show bool) error
	UpdateCache() error
}
