package biz

import (
	"github.com/docker/docker/api/types/network"

	"github.com/TheTNB/panel/internal/http/request"
)

type ContainerNetworkRepo interface {
	List() ([]network.Inspect, error)
	Create(req *request.ContainerNetworkCreate) (string, error)
	Remove(id string) error
	Prune() error
}
