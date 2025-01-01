package biz

import (
	"github.com/tnb-labs/panel/internal/http/request"
	"github.com/tnb-labs/panel/pkg/types"
)

type ContainerNetworkRepo interface {
	List() ([]types.ContainerNetwork, error)
	Create(req *request.ContainerNetworkCreate) (string, error)
	Remove(id string) error
	Prune() error
}
