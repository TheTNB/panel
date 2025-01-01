package biz

import (
	"github.com/tnb-labs/panel/internal/http/request"
	"github.com/tnb-labs/panel/pkg/types"
)

type ContainerImageRepo interface {
	List() ([]types.ContainerImage, error)
	Pull(req *request.ContainerImagePull) error
	Remove(id string) error
	Prune() error
}
