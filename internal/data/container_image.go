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

type containerImageRepo struct {
	cmd string
}

func NewContainerImageRepo(cmd ...string) biz.ContainerImageRepo {
	if len(cmd) == 0 {
		cmd = append(cmd, "docker")
	}
	return &containerImageRepo{
		cmd: cmd[0],
	}
}

// List 列出镜像
func (r *containerImageRepo) List() ([]types.ContainerImage, error) {
	output, err := shell.ExecfWithTimeout(10*time.Second, "%s images -a --format json", r.cmd)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(output, "\n")

	var images []types.ContainerImage
	for _, line := range lines {
		if line == "" {
			continue
		}

		var item struct {
			ID          string `json:"ID"`
			Containers  string `json:"Containers"`
			Repository  string `json:"Repository"`
			Tag         string `json:"Tag"`
			Digest      string `json:"Digest"`
			CreatedAt   string `json:"CreatedAt"`
			Size        string `json:"Size"`
			SharedSize  string `json:"SharedSize"`
			VirtualSize string `json:"VirtualSize"`
		}
		if err = json.Unmarshal([]byte(line), &item); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}

		createdAt, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", item.CreatedAt)
		images = append(images, types.ContainerImage{
			ID:         item.ID,
			Containers: cast.ToInt64(item.Containers),
			Tag:        item.Tag,
			CreatedAt:  createdAt,
			Size:       item.Size,
		})
	}

	return images, nil
}

// Pull 拉取镜像
func (r *containerImageRepo) Pull(req *request.ContainerImagePull) error {
	var sb strings.Builder

	if req.Auth {
		sb.WriteString(fmt.Sprintf("%s login -u %s -p %s", r.cmd, req.Username, req.Password))
		if _, err := shell.ExecfWithTimeout(1*time.Minute, sb.String()); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
		sb.Reset()
	}

	sb.WriteString(fmt.Sprintf("%s pull %s", r.cmd, req.Name))

	if _, err := shell.ExecfWithTimeout(20*time.Minute, sb.String()); err != nil { // nolint: govet
		return fmt.Errorf("pull failed: %w", err)
	}

	return nil
}

// Remove 删除镜像
func (r *containerImageRepo) Remove(id string) error {
	_, err := shell.ExecfWithTimeout(30*time.Second, "%s rmi %s", r.cmd, id)
	return err
}

// Prune 清理未使用的镜像
func (r *containerImageRepo) Prune() error {
	_, err := shell.ExecfWithTimeout(30*time.Second, "%s image prune -f", r.cmd)
	return err
}
