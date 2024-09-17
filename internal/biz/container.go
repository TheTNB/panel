package biz

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	"github.com/TheTNB/panel/internal/http/request"
)

type ContainerRepo interface {
	ListAll() ([]types.Container, error)
	ListByNames(names []string) ([]types.Container, error)
	Create(req *request.ContainerCreate) (string, error)
	Remove(id string) error
	Start(id string) error
	Stop(id string) error
	Restart(id string) error
	Pause(id string) error
	Unpause(id string) error
	Inspect(id string) (types.ContainerJSON, error)
	Kill(id string) error
	Rename(id string, newName string) error
	Stats(id string) (container.StatsResponseReader, error)
	Exist(name string) (bool, error)
	Update(id string, config container.UpdateConfig) error
	Logs(id string) (string, error)
	Prune() error
}
