package services

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type Container struct {
	client *client.Client
}

func NewContainer(sock ...string) *Container {
	if len(sock) == 0 {
		sock[0] = "/run/podman/podman.sock"
	}
	cli, _ := client.NewClientWithOpts(client.WithHost("unix://"+sock[0]), client.WithAPIVersionNegotiation())
	return &Container{
		client: cli,
	}
}

// ContainerListAll 列出所有容器
func (r *Container) ContainerListAll() ([]types.Container, error) {
	containers, err := r.client.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// ContainerListByNames 根据名称列出容器
func (r *Container) ContainerListByNames(names []string) ([]types.Container, error) {
	var options container.ListOptions
	options.All = true
	if len(names) > 0 {
		var array []filters.KeyValuePair
		for _, n := range names {
			array = append(array, filters.Arg("name", n))
		}
		options.Filters = filters.NewArgs(array...)
	}
	containers, err := r.client.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// ContainerCreate 创建容器
func (r *Container) ContainerCreate(name string, config container.Config, host container.HostConfig, network network.NetworkingConfig) (string, error) {
	resp, err := r.client.ContainerCreate(context.Background(), &config, &host, &network, nil, name)
	return resp.ID, err
}

// ContainerRemove 移除容器
func (r *Container) ContainerRemove(id string) error {
	return r.client.ContainerRemove(context.Background(), id, container.RemoveOptions{
		Force: true,
	})
}

// ContainerStart 启动容器
func (r *Container) ContainerStart(id string) error {
	return r.client.ContainerStart(context.Background(), id, container.StartOptions{})
}

// ContainerStop 停止容器
func (r *Container) ContainerStop(id string) error {
	return r.client.ContainerStop(context.Background(), id, container.StopOptions{})
}

// ContainerPause 暂停容器
func (r *Container) ContainerPause(id string) error {
	return r.client.ContainerPause(context.Background(), id)
}

// ContainerUnpause 恢复容器
func (r *Container) ContainerUnpause(id string) error {
	return r.client.ContainerUnpause(context.Background(), id)
}

// ContainerInspect 查看容器
func (r *Container) ContainerInspect(id string) (types.ContainerJSON, error) {
	return r.client.ContainerInspect(context.Background(), id)
}

// ContainerKill 杀死容器
func (r *Container) ContainerKill(id string) error {
	return r.client.ContainerKill(context.Background(), id, "KILL")
}

// ContainerRename 重命名容器
func (r *Container) ContainerRename(id string, newName string) error {
	return r.client.ContainerRename(context.Background(), id, newName)
}

// ContainerStats 查看容器状态
func (r *Container) ContainerStats(id string) (types.ContainerStats, error) {
	return r.client.ContainerStats(context.Background(), id, false)
}

// ContainerExists 判断容器是否存在
func (r *Container) ContainerExists(name string) (bool, error) {
	var options container.ListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	containers, err := r.client.ContainerList(context.Background(), options)
	if err != nil {
		return false, err
	}

	return len(containers) > 0, nil
}

// ContainerUpdate 更新容器
func (r *Container) ContainerUpdate(id string, config container.UpdateConfig) error {
	_, err := r.client.ContainerUpdate(context.Background(), id, config)
	return err
}

// ContainerLogs 查看容器日志
func (r *Container) ContainerLogs(id string) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}
	reader, err := r.client.ContainerLogs(context.Background(), id, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ContainerPrune 清理未使用的容器
func (r *Container) ContainerPrune() error {
	_, err := r.client.ContainersPrune(context.Background(), filters.NewArgs())
	return err
}

// ContainerRestart 重启容器
func (r *Container) ContainerRestart(id string) error {
	return r.client.ContainerRestart(context.Background(), id, container.StopOptions{})
}

// NetworkList 列出网络
func (r *Container) NetworkList() ([]types.NetworkResource, error) {
	return r.client.NetworkList(context.Background(), types.NetworkListOptions{})
}

// NetworkCreate 创建网络
func (r *Container) NetworkCreate(name string) error {
	_, err := r.client.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Driver: "bridge",
	})

	return err
}

// NetworkExist 判断网络是否存在
func (r *Container) NetworkExist(name string) (bool, error) {
	var options types.NetworkListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	networks, err := r.client.NetworkList(context.Background(), options)
	if err != nil {
		return false, err
	}

	return len(networks) > 0, nil
}

// NetworkInspect 查看网络
func (r *Container) NetworkInspect(id string) (types.NetworkResource, error) {
	return r.client.NetworkInspect(context.Background(), id, types.NetworkInspectOptions{})
}

// NetworkConnect 连接网络
func (r *Container) NetworkConnect(id string, containerID string) error {
	return r.client.NetworkConnect(context.Background(), id, containerID, nil)
}

// NetworkDisconnect 断开网络
func (r *Container) NetworkDisconnect(id string, containerID string) error {
	return r.client.NetworkDisconnect(context.Background(), id, containerID, true)
}

// NetworkPrune 清理未使用的网络
func (r *Container) NetworkPrune() error {
	_, err := r.client.NetworksPrune(context.Background(), filters.NewArgs())
	return err
}

// NetworkRemove 删除网络
func (r *Container) NetworkRemove(id string) error {
	return r.client.NetworkRemove(context.Background(), id)
}

// ImageList 列出镜像
func (r *Container) ImageList() ([]image.Summary, error) {
	return r.client.ImageList(context.Background(), types.ImageListOptions{
		All: true,
	})
}

// ImageExist 判断镜像是否存在
func (r *Container) ImageExist(str string) (bool, error) {
	var options types.ImageListOptions
	options.Filters = filters.NewArgs(filters.Arg("reference", str))
	images, err := r.client.ImageList(context.Background(), options)
	if err != nil {
		return false, err
	}

	return len(images) > 0, nil
}

// ImagePull 拉取镜像
func (r *Container) ImagePull(str string) error {
	_, err := r.client.ImagePull(context.Background(), str, types.ImagePullOptions{})
	return err
}

// ImageRemove 删除镜像
func (r *Container) ImageRemove(id string) error {
	_, err := r.client.ImageRemove(context.Background(), id, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	return err
}

// ImagePrune 清理未使用的镜像
func (r *Container) ImagePrune() error {
	_, err := r.client.ImagesPrune(context.Background(), filters.NewArgs())
	return err
}

// ImageInspect 查看镜像
func (r *Container) ImageInspect(id string) (types.ImageInspect, error) {
	img, _, err := r.client.ImageInspectWithRaw(context.Background(), id)
	return img, err
}

// VolumeCreate 创建存储卷
func (r *Container) VolumeCreate(name string, options, labels map[string]string) (volume.Volume, error) {
	return r.client.VolumeCreate(context.Background(), volume.CreateOptions{
		Name:       name,
		Driver:     "local",
		DriverOpts: options,
		Labels:     labels,
	})
}

// VolumeList 列出存储卷
func (r *Container) VolumeList() ([]*volume.Volume, error) {
	volumes, err := r.client.VolumeList(context.Background(), volume.ListOptions{})
	return volumes.Volumes, err
}

// VolumeExist 判断存储卷是否存在
func (r *Container) VolumeExist(name string) (bool, error) {
	var options volume.ListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	volumes, err := r.client.VolumeList(context.Background(), options)
	if err != nil {
		return false, err
	}

	return len(volumes.Volumes) > 0, nil
}

// VolumeInspect 查看存储卷
func (r *Container) VolumeInspect(id string) (volume.Volume, error) {
	return r.client.VolumeInspect(context.Background(), id)
}

// VolumeRemove 删除存储卷
func (r *Container) VolumeRemove(id string) error {
	return r.client.VolumeRemove(context.Background(), id, true)
}

// VolumePrune 清理未使用的存储卷
func (r *Container) VolumePrune() error {
	_, err := r.client.VolumesPrune(context.Background(), filters.NewArgs())
	return err
}
