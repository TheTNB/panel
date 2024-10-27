package data

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/types"
)

type containerVolumeRepo struct {
	client *client.Client
}

func NewContainerVolumeRepo(sock ...string) biz.ContainerVolumeRepo {
	if len(sock) == 0 {
		sock = append(sock, "/run/podman/podman.sock")
	}
	cli, _ := client.NewClientWithOpts(client.WithHost("unix://"+sock[0]), client.WithAPIVersionNegotiation())
	return &containerVolumeRepo{
		client: cli,
	}
}

// List 列出存储卷
func (r *containerVolumeRepo) List() ([]*volume.Volume, error) {
	volumes, err := r.client.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		return nil, err
	}
	return volumes.Volumes, err
}

// Create 创建存储卷
func (r *containerVolumeRepo) Create(req *request.ContainerVolumeCreate) (volume.Volume, error) {
	return r.client.VolumeCreate(context.Background(), volume.CreateOptions{
		Name:       req.Name,
		Driver:     req.Driver,
		DriverOpts: types.KVToMap(req.Options),
		Labels:     types.KVToMap(req.Labels),
	})
}

// Remove 删除存储卷
func (r *containerVolumeRepo) Remove(id string) error {
	return r.client.VolumeRemove(context.Background(), id, true)
}

// Prune 清理未使用的存储卷
func (r *containerVolumeRepo) Prune() error {
	_, err := r.client.VolumesPrune(context.Background(), filters.NewArgs())
	return err
}
