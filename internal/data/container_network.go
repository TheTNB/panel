package data

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/types"
)

type containerNetworkRepo struct {
	cmd string
}

func NewContainerNetworkRepo(cmd ...string) biz.ContainerNetworkRepo {
	if len(cmd) == 0 {
		cmd = append(cmd, "docker")
	}
	return &containerNetworkRepo{
		cmd: cmd[0],
	}
}

// List 列出网络
func (r *containerNetworkRepo) List() ([]types.ContainerNetwork, error) {
	output, err := shell.ExecfWithTimeout(10*time.Second, "%s network ls --format json", r.cmd)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(output, "\n")

	var networks []types.ContainerNetwork
	for _, line := range lines {
		if line == "" {
			continue
		}

		var item struct {
			CreatedAt string `json:"CreatedAt"`
			Driver    string `json:"Driver"`
			ID        string `json:"ID"`
			IPv6      string `json:"IPv6"`
			Internal  string `json:"Internal"`
			Labels    string `json:"Labels"`
			Name      string `json:"Name"`
			Scope     string `json:"Scope"`
		}
		if err = json.Unmarshal([]byte(line), &item); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}

		output, err = shell.ExecfWithTimeout(10*time.Second, "%s network inspect %s", r.cmd, item.ID)
		if err != nil {
			return nil, fmt.Errorf("inspect failed: %w", err)
		}
		var inspect []struct {
			Name       string    `json:"Name"`
			Id         string    `json:"Id"`
			Created    time.Time `json:"Created"`
			Scope      string    `json:"Scope"`
			Driver     string    `json:"Driver"`
			EnableIPv6 bool      `json:"EnableIPv6"`
			IPAM       struct {
				Driver  string            `json:"Driver"`
				Options map[string]string `json:"Options"`
				Config  []struct {
					Subnet     string            `json:"Subnet"`
					IPRange    string            `json:"IPRange"`
					Gateway    string            `json:"Gateway"`
					AuxAddress map[string]string `json:"AuxiliaryAddresses"`
				} `json:"Config"`
			} `json:"IPAM"`
			Internal   bool              `json:"Internal"`
			Attachable bool              `json:"Attachable"`
			Ingress    bool              `json:"Ingress"`
			ConfigOnly bool              `json:"ConfigOnly"`
			Options    map[string]string `json:"Options"`
			Labels     map[string]string `json:"Labels"`
		}
		if err = json.Unmarshal([]byte(output), &inspect); err != nil {
			return nil, fmt.Errorf("unmarshal inspect failed: %w", err)
		}
		if len(inspect) == 0 {
			return nil, fmt.Errorf("inspect empty")
		}

		var ipamConfigs []types.ContainerNetworkIPAMConfig
		for _, ipam := range inspect[0].IPAM.Config {
			ipamConfigs = append(ipamConfigs, types.ContainerNetworkIPAMConfig{
				Subnet:     ipam.Subnet,
				IPRange:    ipam.IPRange,
				Gateway:    ipam.Gateway,
				AuxAddress: ipam.AuxAddress,
			})
		}

		createdAt, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", item.CreatedAt)
		networks = append(networks, types.ContainerNetwork{
			ID:         item.ID,
			Name:       item.Name,
			Driver:     item.Driver,
			IPv6:       cast.ToBool(item.IPv6),
			Internal:   cast.ToBool(item.Internal),
			Attachable: cast.ToBool(inspect[0].Attachable),
			Ingress:    cast.ToBool(inspect[0].Ingress),
			Scope:      item.Scope,
			CreatedAt:  createdAt,
			IPAM: types.ContainerNetworkIPAM{
				Driver:  inspect[0].IPAM.Driver,
				Options: types.MapToKV(inspect[0].IPAM.Options),
				Config:  ipamConfigs,
			},
			Options: types.MapToKV(inspect[0].Options),
			Labels:  types.SliceToKV(strings.Split(item.Labels, ",")),
		})
	}

	return networks, nil
}

// Create 创建网络
func (r *containerNetworkRepo) Create(req *request.ContainerNetworkCreate) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s network create --driver %s", r.cmd, req.Driver))
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

	return shell.ExecfWithTimeout(10*time.Second, "%s", sb.String()) // nolint: govet
}

// Remove 删除网络
func (r *containerNetworkRepo) Remove(id string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s network rm -f %s", r.cmd, id)
	return err
}

// Prune 清理未使用的网络
func (r *containerNetworkRepo) Prune() error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s network prune -f", r.cmd)
	return err
}
