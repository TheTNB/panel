package data

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/types"
)

type containerVolumeRepo struct {
	cmd string
}

func NewContainerVolumeRepo(cmd ...string) biz.ContainerVolumeRepo {
	if len(cmd) == 0 {
		cmd = append(cmd, "docker")
	}
	return &containerVolumeRepo{
		cmd: cmd[0],
	}
}

// List 列出存储卷
func (r *containerVolumeRepo) List() ([]types.ContainerVolume, error) {
	output, err := shell.ExecfWithTimeout(10*time.Second, "%s volume ls --format json", r.cmd)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(output, "\n")

	var volumes []types.ContainerVolume
	for _, line := range lines {
		if line == "" {
			continue
		}

		var item struct {
			Availability string `json:"Availability"`
			Driver       string `json:"Driver"`
			Group        string `json:"Group"`
			Labels       string `json:"Labels"`
			Links        string `json:"Links"`
			Mountpoint   string `json:"Mountpoint"`
			Name         string `json:"Name"`
			Scope        string `json:"Scope"`
			Size         string `json:"Size"`
			Status       string `json:"Status"`
		}
		if err = json.Unmarshal([]byte(line), &item); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}

		output, err = shell.ExecfWithTimeout(10*time.Second, "%s volume inspect %s", r.cmd, item.Name)
		if err != nil {
			return nil, fmt.Errorf("inspect failed: %w", err)
		}
		var inspect []struct {
			CreatedAt  time.Time         `json:"CreatedAt"`
			Driver     string            `json:"Driver"`
			Labels     map[string]string `json:"Labels"`
			Mountpoint string            `json:"Mountpoint"`
			Name       string            `json:"Name"`
			Options    map[string]string `json:"Options"`
			Scope      string            `json:"Scope"`
		}
		if err = json.Unmarshal([]byte(output), &inspect); err != nil {
			return nil, fmt.Errorf("unmarshal inspect failed: %w", err)
		}
		if len(inspect) == 0 {
			return nil, fmt.Errorf("inspect empty")
		}

		volumes = append(volumes, types.ContainerVolume{
			Name:       item.Name,
			Driver:     item.Driver,
			Scope:      item.Scope,
			MountPoint: item.Mountpoint,
			CreatedAt:  inspect[0].CreatedAt,
			Options:    types.MapToKV(inspect[0].Options),
			Labels:     types.SliceToKV(strings.Split(item.Labels, ",")),
		})
	}

	return volumes, nil
}

// Create 创建存储卷
func (r *containerVolumeRepo) Create(req *request.ContainerVolumeCreate) (string, error) {
	return "", nil
}

// Remove 删除存储卷
func (r *containerVolumeRepo) Remove(id string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s volume rm -f %s", r.cmd, id)
	return err
}

// Prune 清理未使用的存储卷
func (r *containerVolumeRepo) Prune() error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s volume prune -f", r.cmd)
	return err
}
