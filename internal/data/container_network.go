package data

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/tnb-labs/panel/internal/biz"
	"github.com/tnb-labs/panel/internal/http/request"
	"github.com/tnb-labs/panel/pkg/shell"
	"github.com/tnb-labs/panel/pkg/types"
	"github.com/tnb-labs/panel/pkg/types/docker/network"
)

type containerNetworkRepo struct {
	client *resty.Client
}

func NewContainerNetworkRepo() biz.ContainerNetworkRepo {
	return &containerNetworkRepo{
		client: getDockerClient("/var/run/docker.sock"),
	}
}

// List 列出网络
func (r *containerNetworkRepo) List() ([]types.ContainerNetwork, error) {
	var resp []network.Network
	_, err := r.client.R().SetResult(&resp).Get("/networks")
	if err != nil {
		return nil, err
	}

	var networks []types.ContainerNetwork
	for _, item := range resp {
		ipamConfigs := make([]types.ContainerNetworkIPAMConfig, 0)
		for _, ipam := range item.IPAM.Config {
			ipamConfigs = append(ipamConfigs, types.ContainerNetworkIPAMConfig{
				Subnet:     ipam.Subnet,
				IPRange:    ipam.IPRange,
				Gateway:    ipam.Gateway,
				AuxAddress: ipam.AuxAddress,
			})
		}
		networks = append(networks, types.ContainerNetwork{
			ID:         item.ID,
			Name:       item.Name,
			Driver:     item.Driver,
			IPv6:       item.EnableIPv6,
			Internal:   item.Internal,
			Attachable: item.Attachable,
			Ingress:    item.Ingress,
			Scope:      item.Scope,
			CreatedAt:  item.Created,
			IPAM: types.ContainerNetworkIPAM{
				Driver:  item.IPAM.Driver,
				Options: types.MapToKV(item.IPAM.Options),
				Config:  ipamConfigs,
			},
			Options: types.MapToKV(item.Options),
			Labels:  types.MapToKV(item.Labels),
		})
	}

	slices.SortFunc(networks, func(a types.ContainerNetwork, b types.ContainerNetwork) int {
		return strings.Compare(a.Name, b.Name)
	})

	return networks, nil
}

// Create 创建网络
func (r *containerNetworkRepo) Create(req *request.ContainerNetworkCreate) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("docker network create --driver %s", req.Driver))
	sb.WriteString(fmt.Sprintf(" %s", req.Name))

	if req.Ipv4.Enabled {
		sb.WriteString(fmt.Sprintf(" --subnet %s", req.Ipv4.Subnet))
		sb.WriteString(fmt.Sprintf(" --gateway %s", req.Ipv4.Gateway))
		if req.Ipv4.IPRange != "" {
			sb.WriteString(fmt.Sprintf(" --ip-range %s", req.Ipv4.IPRange))
		}
	}
	if req.Ipv6.Enabled {
		sb.WriteString(fmt.Sprintf(" --subnet %s", req.Ipv6.Subnet))
		sb.WriteString(fmt.Sprintf(" --gateway %s", req.Ipv6.Gateway))
		if req.Ipv6.IPRange != "" {
			sb.WriteString(fmt.Sprintf(" --ip-range %s", req.Ipv6.IPRange))
		}
	}
	for _, label := range req.Labels {
		sb.WriteString(fmt.Sprintf(" --label %s=%s", label.Key, label.Value))
	}
	for _, option := range req.Options {
		sb.WriteString(fmt.Sprintf(" --opt %s=%s", option.Key, option.Value))
	}

	return shell.Exec(sb.String())
}

// Remove 删除网络
func (r *containerNetworkRepo) Remove(id string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker network rm -f %s", id)
	return err
}

// Prune 清理未使用的网络
func (r *containerNetworkRepo) Prune() error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker network prune -f")
	return err
}
