package biz

import (
	"github.com/tnb-labs/panel/internal/http/request"
	"github.com/tnb-labs/panel/pkg/types"
)

type ContainerVolumeRepo interface {
	List() ([]types.ContainerVolume, error)
	Create(req *request.ContainerVolumeCreate) (string, error)
	Remove(id string) error
	Prune() error
}
