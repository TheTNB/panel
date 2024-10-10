package biz

import (
	"time"

	"github.com/TheTNB/panel/pkg/api"
)

type App struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Slug      string    `gorm:"not null;unique" json:"slug"`
	Channel   string    `gorm:"not null" json:"channel"`
	Version   string    `gorm:"not null" json:"version"`
	Show      bool      `gorm:"not null" json:"show"`
	ShowOrder int       `gorm:"not null" json:"show_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppRepo interface {
	All() api.Apps
	Get(slug string) (*api.App, error)
	UpdateExist(slug string) bool
	Installed() ([]*App, error)
	GetInstalled(slug string) (*App, error)
	GetInstalledAll(query string, cond ...string) ([]*App, error)
	GetHomeShow() ([]map[string]string, error)
	IsInstalled(query string, cond ...string) (bool, error)
	Install(channel, slug string) error
	UnInstall(slug string) error
	Update(slug string) error
	UpdateShow(slug string, show bool) error
	UpdateCache() error
}
