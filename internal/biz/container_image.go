package biz

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"

	"github.com/TheTNB/panel/internal/http/request"
)

type ContainerImageRepo interface {
	List() ([]image.Summary, error)
	Exist(id string) (bool, error)
	Pull(req request.ContainerImagePull) error
	Remove(id string) error
	Prune() error
	Inspect(id string) (types.ImageInspect, error)
}
