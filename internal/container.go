package internal

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
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
	ContainerStats(id string) (types.ContainerStats, error)
	ContainerExist(name string) (bool, error)
	ContainerUpdate(id string, config container.UpdateConfig) error
	ContainerLogs(id string) (string, error)
	ContainerPrune() error
	NetworkList() ([]types.NetworkResource, error)
	NetworkCreate(name string) error
	NetworkRemove(id string) error
	NetworkExist(name string) (bool, error)
	NetworkInspect(id string) (types.NetworkResource, error)
	NetworkConnect(networkID string, containerID string) error
	NetworkDisconnect(networkID string, containerID string) error
	NetworkPrune() error
	ImageList() ([]image.Summary, error)
	ImageExist(reference string) (bool, error)
	ImagePull(reference string) error
	ImageRemove(imageID string) error
	ImagePrune() error
	ImageInspect(imageID string) (types.ImageInspect, error)
	VolumeList() ([]*volume.Volume, error)
	VolumeCreate(name string, options, labels map[string]string) (volume.Volume, error)
	VolumeExist(name string) (bool, error)
	VolumeInspect(volumeID string) (volume.Volume, error)
	VolumeRemove(volumeID string) error
	VolumePrune() error
	SliceToMap(slice []string) map[string]string
}
