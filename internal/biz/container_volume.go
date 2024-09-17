package biz

import (
	"github.com/docker/docker/api/types/volume"

	"github.com/TheTNB/panel/internal/http/request"
)

type ContainerVolumeRepo interface {
	List() ([]*volume.Volume, error)
	Create(req request.ContainerVolumeCreate) (volume.Volume, error)
	Exist(name string) (bool, error)
	Inspect(id string) (volume.Volume, error)
	Remove(id string) error
	Prune() error
}
