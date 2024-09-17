package data

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	paneltypes "github.com/TheTNB/panel/pkg/types"
)

type containerRepo struct {
	client *client.Client
}

func NewContainerRepo(sock ...string) biz.ContainerRepo {
	if len(sock) == 0 {
		sock = append(sock, "/run/podman/podman.sock")
	}
	cli, _ := client.NewClientWithOpts(client.WithHost("unix://"+sock[0]), client.WithAPIVersionNegotiation())
	return &containerRepo{
		client: cli,
	}
}

// ListAll 列出所有容器
func (r *containerRepo) ListAll() ([]types.Container, error) {
	containers, err := r.client.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// ListByNames 根据名称列出容器
func (r *containerRepo) ListByNames(names []string) ([]types.Container, error) {
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

// Create 创建容器
func (r *containerRepo) Create(req *request.ContainerCreate) (string, error) {
	var hostConf container.HostConfig
	var networkConf network.NetworkingConfig

	portMap := make(nat.PortMap)
	for _, port := range req.Ports {
		if port.ContainerStart-port.ContainerEnd != port.HostStart-port.HostEnd {
			return "", fmt.Errorf("容器端口和主机端口数量不匹配（容器: %d 主机: %d）", port.ContainerStart-port.ContainerEnd, port.HostStart-port.HostEnd)
		}
		if port.ContainerStart > port.ContainerEnd || port.HostStart > port.HostEnd || port.ContainerStart < 1 || port.HostStart < 1 {
			return "", fmt.Errorf("端口范围不正确")
		}

		count := 0
		for host := port.HostStart; host <= port.HostEnd; host++ {
			bindItem := nat.PortBinding{HostPort: strconv.Itoa(host), HostIP: port.Host}
			portMap[nat.Port(fmt.Sprintf("%d/%s", port.ContainerStart+count, port.Protocol))] = []nat.PortBinding{bindItem}
			count++
		}
	}

	exposed := make(nat.PortSet)
	for port := range portMap {
		exposed[port] = struct{}{}
	}

	if req.Network != "" {
		switch req.Network {
		case "host", "none", "bridge":
			hostConf.NetworkMode = container.NetworkMode(req.Network)
		}
		networkConf.EndpointsConfig = map[string]*network.EndpointSettings{req.Network: {}}
	} else {
		networkConf = network.NetworkingConfig{}
	}

	hostConf.Privileged = req.Privileged
	hostConf.AutoRemove = req.AutoRemove
	hostConf.CPUShares = req.CPUShares
	hostConf.PublishAllPorts = req.PublishAllPorts
	hostConf.RestartPolicy = container.RestartPolicy{Name: container.RestartPolicyMode(req.RestartPolicy)}
	if req.RestartPolicy == "on-failure" {
		hostConf.RestartPolicy.MaximumRetryCount = 5
	}
	hostConf.NanoCPUs = req.CPUs * 1000000000
	hostConf.Memory = req.Memory * 1024 * 1024
	hostConf.MemorySwap = 0
	hostConf.PortBindings = portMap
	hostConf.Binds = []string{}

	volumes := make(map[string]struct{})
	for _, v := range req.Volumes {
		volumes[v.Container] = struct{}{}
		hostConf.Binds = append(hostConf.Binds, fmt.Sprintf("%s:%s:%s", v.Host, v.Container, v.Mode))
	}

	resp, err := r.client.ContainerCreate(context.Background(), &container.Config{
		Image:        req.Image,
		Env:          paneltypes.KVToSlice(req.Env),
		Entrypoint:   req.Entrypoint,
		Cmd:          req.Command,
		Labels:       paneltypes.KVToMap(req.Labels),
		ExposedPorts: exposed,
		OpenStdin:    req.OpenStdin,
		Tty:          req.Tty,
		Volumes:      volumes,
	}, &hostConf, &networkConf, nil, req.Name)
	if err != nil {
		return "", err
	}

	return resp.ID, err
}

// Remove 移除容器
func (r *containerRepo) Remove(id string) error {
	return r.client.ContainerRemove(context.Background(), id, container.RemoveOptions{
		Force: true,
	})
}

// Start 启动容器
func (r *containerRepo) Start(id string) error {
	return r.client.ContainerStart(context.Background(), id, container.StartOptions{})
}

// Stop 停止容器
func (r *containerRepo) Stop(id string) error {
	return r.client.ContainerStop(context.Background(), id, container.StopOptions{})
}

// Restart 重启容器
func (r *containerRepo) Restart(id string) error {
	return r.client.ContainerRestart(context.Background(), id, container.StopOptions{})
}

// Pause 暂停容器
func (r *containerRepo) Pause(id string) error {
	return r.client.ContainerPause(context.Background(), id)
}

// Unpause 恢复容器
func (r *containerRepo) Unpause(id string) error {
	return r.client.ContainerUnpause(context.Background(), id)
}

// Inspect 查看容器
func (r *containerRepo) Inspect(id string) (types.ContainerJSON, error) {
	return r.client.ContainerInspect(context.Background(), id)
}

// Kill 杀死容器
func (r *containerRepo) Kill(id string) error {
	return r.client.ContainerKill(context.Background(), id, "KILL")
}

// Rename 重命名容器
func (r *containerRepo) Rename(id string, newName string) error {
	return r.client.ContainerRename(context.Background(), id, newName)
}

// Stats 查看容器状态
func (r *containerRepo) Stats(id string) (container.StatsResponseReader, error) {
	return r.client.ContainerStats(context.Background(), id, false)
}

// Exist 判断容器是否存在
func (r *containerRepo) Exist(name string) (bool, error) {
	var options container.ListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	containers, err := r.client.ContainerList(context.Background(), options)
	if err != nil {
		return false, err
	}

	return len(containers) > 0, nil
}

// Update 更新容器
func (r *containerRepo) Update(id string, config container.UpdateConfig) error {
	_, err := r.client.ContainerUpdate(context.Background(), id, config)
	return err
}

// Logs 查看容器日志
func (r *containerRepo) Logs(id string) (string, error) {
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

// Prune 清理未使用的容器
func (r *containerRepo) Prune() error {
	_, err := r.client.ContainersPrune(context.Background(), filters.NewArgs())
	return err
}
