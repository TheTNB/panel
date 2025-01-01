package biz

import (
	"github.com/tnb-labs/panel/internal/http/request"
	"github.com/tnb-labs/panel/pkg/types"
)

type ContainerRepo interface {
	ListAll() ([]types.Container, error)
	ListByName(name string) ([]types.Container, error)
	Create(req *request.ContainerCreate) (string, error)
	Remove(id string) error
	Start(id string) error
	Stop(id string) error
	Restart(id string) error
	Pause(id string) error
	Unpause(id string) error
	Kill(id string) error
	Rename(id string, newName string) error
	Logs(id string) (string, error)
	Prune() error
}
