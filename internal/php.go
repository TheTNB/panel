package internal

import (
	"github.com/TheTNB/panel/v2/pkg/types"
)

type PHP interface {
	Status() (bool, error)
	Reload() error
	Start() error
	Stop() error
	Restart() error
	GetConfig() (string, error)
	SaveConfig(config string) error
	GetFPMConfig() (string, error)
	SaveFPMConfig(config string) error
	Load() ([]types.NV, error)
	GetErrorLog() (string, error)
	GetSlowLog() (string, error)
	ClearErrorLog() error
	ClearSlowLog() error
	GetExtensions() ([]types.PHPExtension, error)
	InstallExtension(slug string) error
	UninstallExtension(slug string) error
}
