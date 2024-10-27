package biz

import (
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/types"
)

type ContainerImageRepo interface {
	List() ([]types.ContainerImage, error)
	Pull(req *request.ContainerImagePull) error
	Remove(id string) error
	Prune() error
}
