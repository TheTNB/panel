package services

import (
	"context"
	"encoding/base64"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/goravel/framework/support/json"

	requests "github.com/TheTNB/panel/v2/app/http/requests/container"
	paneltypes "github.com/TheTNB/panel/v2/pkg/types"
)

type Container struct {
	client *client.Client
}

func NewContainer(sock ...string) Container {
	if len(sock) == 0 {
		sock = append(sock, "/run/podman/podman.sock")
	}
	cli, _ := client.NewClientWithOpts(client.WithHost("unix://"+sock[0]), client.WithAPIVersionNegotiation())
	return Container{
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

// ContainerRestart 重启容器
func (r *Container) ContainerRestart(id string) error {
	return r.client.ContainerRestart(context.Background(), id, container.StopOptions{})
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
func (r *Container) ContainerStats(id string) (container.StatsResponseReader, error) {
	return r.client.ContainerStats(context.Background(), id, false)
}

// ContainerExist 判断容器是否存在
func (r *Container) ContainerExist(name string) (bool, error) {
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

// NetworkList 列出网络
func (r *Container) NetworkList() ([]network.Inspect, error) {
	return r.client.NetworkList(context.Background(), network.ListOptions{})
}

// NetworkCreate 创建网络
func (r *Container) NetworkCreate(config requests.NetworkCreate) (string, error) {
	var ipamConfigs []network.IPAMConfig
	if config.Ipv4.Enabled {
		ipamConfigs = append(ipamConfigs, network.IPAMConfig{
			Subnet:  config.Ipv4.Subnet,
			Gateway: config.Ipv4.Gateway,
			IPRange: config.Ipv4.IPRange,
		})
	}
	if config.Ipv6.Enabled {
		ipamConfigs = append(ipamConfigs, network.IPAMConfig{
			Subnet:  config.Ipv6.Subnet,
			Gateway: config.Ipv6.Gateway,
			IPRange: config.Ipv6.IPRange,
		})
	}

	options := network.CreateOptions{
		EnableIPv6: &config.Ipv6.Enabled,
		Driver:     config.Driver,
		Options:    r.KVToMap(config.Options),
		Labels:     r.KVToMap(config.Labels),
	}
	if len(ipamConfigs) > 0 {
		options.IPAM = &network.IPAM{
			Config: ipamConfigs,
		}
	}

	resp, err := r.client.NetworkCreate(context.Background(), config.Name, options)
	return resp.ID, err
}

// NetworkRemove 删除网络
func (r *Container) NetworkRemove(id string) error {
	return r.client.NetworkRemove(context.Background(), id)
}

// NetworkExist 判断网络是否存在
func (r *Container) NetworkExist(name string) (bool, error) {
	var options network.ListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	networks, err := r.client.NetworkList(context.Background(), options)
	if err != nil {
		return false, err
	}

	return len(networks) > 0, nil
}

// NetworkInspect 查看网络
func (r *Container) NetworkInspect(id string) (network.Inspect, error) {
	return r.client.NetworkInspect(context.Background(), id, network.InspectOptions{})
}

// NetworkConnect 连接网络
func (r *Container) NetworkConnect(networkID string, containerID string) error {
	return r.client.NetworkConnect(context.Background(), networkID, containerID, nil)
}

// NetworkDisconnect 断开网络
func (r *Container) NetworkDisconnect(networkID string, containerID string) error {
	return r.client.NetworkDisconnect(context.Background(), networkID, containerID, true)
}

// NetworkPrune 清理未使用的网络
func (r *Container) NetworkPrune() error {
	_, err := r.client.NetworksPrune(context.Background(), filters.NewArgs())
	return err
}

// ImageList 列出镜像
func (r *Container) ImageList() ([]image.Summary, error) {
	return r.client.ImageList(context.Background(), image.ListOptions{
		All: true,
	})
}

// ImageExist 判断镜像是否存在
func (r *Container) ImageExist(id string) (bool, error) {
	var options image.ListOptions
	options.Filters = filters.NewArgs(filters.Arg("reference", id))
	images, err := r.client.ImageList(context.Background(), options)
	if err != nil {
		return false, err
	}

	return len(images) > 0, nil
}

// ImagePull 拉取镜像
func (r *Container) ImagePull(config requests.ImagePull) error {
	options := image.PullOptions{}
	if config.Auth {
		authConfig := registry.AuthConfig{
			Username: config.Username,
			Password: config.Password,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		options.RegistryAuth = authStr
	}

	out, err := r.client.ImagePull(context.Background(), config.Name, options)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(io.Discard, out)
	return err
}

// ImageRemove 删除镜像
func (r *Container) ImageRemove(id string) error {
	_, err := r.client.ImageRemove(context.Background(), id, image.RemoveOptions{
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

// VolumeList 列出存储卷
func (r *Container) VolumeList() ([]*volume.Volume, error) {
	volumes, err := r.client.VolumeList(context.Background(), volume.ListOptions{})
	return volumes.Volumes, err
}

// VolumeCreate 创建存储卷
func (r *Container) VolumeCreate(config requests.VolumeCreate) (volume.Volume, error) {
	return r.client.VolumeCreate(context.Background(), volume.CreateOptions{
		Name:       config.Name,
		Driver:     config.Driver,
		DriverOpts: r.KVToMap(config.Options),
		Labels:     r.KVToMap(config.Labels),
	})
}

// VolumeExist 判断存储卷是否存在
func (r *Container) VolumeExist(id string) (bool, error) {
	var options volume.ListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", id))
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

// KVToMap 将 key-value 切片转换为 map
func (r *Container) KVToMap(kvs []paneltypes.KV) map[string]string {
	m := make(map[string]string)
	for _, item := range kvs {
		m[item.Key] = item.Value
	}

	return m
}

// KVToSlice 将 key-value 切片转换为 key=value 切片
func (r *Container) KVToSlice(kvs []paneltypes.KV) []string {
	var s []string
	for _, item := range kvs {
		s = append(s, item.Key+"="+item.Value)
	}

	return s
}
