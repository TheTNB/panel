package internal

import (
	"github.com/TheTNB/panel/app/models"
	"github.com/TheTNB/panel/types"
)

type Plugin interface {
	AllInstalled() ([]models.Plugin, error)
	All() []types.Plugin
	GetBySlug(slug string) types.Plugin
	GetInstalledBySlug(slug string) models.Plugin
	Install(slug string) error
	Uninstall(slug string) error
	Update(slug string) error
}
