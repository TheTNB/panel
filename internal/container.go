package internal

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"

	requests "github.com/TheTNB/panel/v2/app/http/requests/container"
	paneltypes "github.com/TheTNB/panel/v2/pkg/types"
)

type Container interface {
	ContainerListAll() ([]types.Container, error)
	ContainerListByNames(names []string) ([]types.Container, error)
	ContainerCreate(name string, config container.Config, host container.HostConfig, networkConfig network.NetworkingConfig) (string, error)
	ContainerRemove(id string) error
	ContainerStart(id string) error
	ContainerStop(id string) error
	ContainerRestart(id string) error
	ContainerPause(id string) error
	ContainerUnpause(id string) error
	ContainerInspect(id string) (types.ContainerJSON, error)
	ContainerKill(id string) error
	ContainerRename(id string, newName string) error
	ContainerStats(id string) (container.StatsResponseReader, error)
	ContainerExist(name string) (bool, error)
	ContainerUpdate(id string, config container.UpdateConfig) error
	ContainerLogs(id string) (string, error)
	ContainerPrune() error
	NetworkList() ([]network.Inspect, error)
	NetworkCreate(config requests.NetworkCreate) (string, error)
	NetworkRemove(id string) error
	NetworkExist(name string) (bool, error)
	NetworkInspect(id string) (network.Inspect, error)
	NetworkConnect(networkID string, containerID string) error
	NetworkDisconnect(networkID string, containerID string) error
	NetworkPrune() error
	ImageList() ([]image.Summary, error)
	ImageExist(id string) (bool, error)
	ImagePull(config requests.ImagePull) error
	ImageRemove(id string) error
	ImagePrune() error
	ImageInspect(id string) (types.ImageInspect, error)
	VolumeList() ([]*volume.Volume, error)
	VolumeCreate(config requests.VolumeCreate) (volume.Volume, error)
	VolumeExist(name string) (bool, error)
	VolumeInspect(id string) (volume.Volume, error)
	VolumeRemove(id string) error
	VolumePrune() error
	KVToMap(kvs []paneltypes.KV) map[string]string
	KVToSlice(kvs []paneltypes.KV) []string
}
