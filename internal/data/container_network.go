package data

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/types"
)

type containerNetworkRepo struct {
	client *client.Client
}

func NewContainerNetworkRepo(sock ...string) biz.ContainerNetworkRepo {
	if len(sock) == 0 {
		sock = append(sock, "/run/podman/podman.sock")
	}
	cli, _ := client.NewClientWithOpts(client.WithHost("unix://"+sock[0]), client.WithAPIVersionNegotiation())
	return &containerNetworkRepo{
		client: cli,
	}
}

// List 列出网络
func (r *containerNetworkRepo) List() ([]network.Inspect, error) {
	return r.client.NetworkList(context.Background(), network.ListOptions{})
}

// Create 创建网络
func (r *containerNetworkRepo) Create(req request.ContainerNetworkCreate) (string, error) {
	var ipamConfigs []network.IPAMConfig
	if req.Ipv4.Enabled {
		ipamConfigs = append(ipamConfigs, network.IPAMConfig{
			Subnet:  req.Ipv4.Subnet,
			Gateway: req.Ipv4.Gateway,
			IPRange: req.Ipv4.IPRange,
		})
	}
	if req.Ipv6.Enabled {
		ipamConfigs = append(ipamConfigs, network.IPAMConfig{
			Subnet:  req.Ipv6.Subnet,
			Gateway: req.Ipv6.Gateway,
			IPRange: req.Ipv6.IPRange,
		})
	}

	options := network.CreateOptions{
		EnableIPv6: &req.Ipv6.Enabled,
		Driver:     req.Driver,
		Options:    types.KVToMap(req.Options),
		Labels:     types.KVToMap(req.Labels),
	}
	if len(ipamConfigs) > 0 {
		options.IPAM = &network.IPAM{
			Config: ipamConfigs,
		}
	}

	resp, err := r.client.NetworkCreate(context.Background(), req.Name, options)
	return resp.ID, err
}

// Remove 删除网络
func (r *containerNetworkRepo) Remove(id string) error {
	return r.client.NetworkRemove(context.Background(), id)
}

// Exist 判断网络是否存在
func (r *containerNetworkRepo) Exist(name string) (bool, error) {
	var options network.ListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	networks, err := r.client.NetworkList(context.Background(), options)
	if err != nil {
		return false, err
	}

	return len(networks) > 0, nil
}

// Inspect 查看网络
func (r *containerNetworkRepo) Inspect(id string) (network.Inspect, error) {
	return r.client.NetworkInspect(context.Background(), id, network.InspectOptions{})
}

// Connect 连接网络
func (r *containerNetworkRepo) Connect(networkID string, containerID string) error {
	return r.client.NetworkConnect(context.Background(), networkID, containerID, nil)
}

// Disconnect 断开网络
func (r *containerNetworkRepo) Disconnect(networkID string, containerID string) error {
	return r.client.NetworkDisconnect(context.Background(), networkID, containerID, true)
}

// Prune 清理未使用的网络
func (r *containerNetworkRepo) Prune() error {
	_, err := r.client.NetworksPrune(context.Background(), filters.NewArgs())
	return err
}
