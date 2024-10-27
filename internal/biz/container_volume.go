package biz

import (
	"github.com/docker/docker/api/types/volume"

	"github.com/TheTNB/panel/internal/http/request"
)

type ContainerVolumeRepo interface {
	List() ([]*volume.Volume, error)
	Create(req *request.ContainerVolumeCreate) (volume.Volume, error)
	Remove(id string) error
	Prune() error
}
