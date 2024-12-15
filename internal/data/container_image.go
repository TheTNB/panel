package data

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/pkg/types"
	"github.com/TheTNB/panel/pkg/types/docker/image"
)

type containerImageRepo struct {
	client *resty.Client
}

func NewContainerImageRepo() biz.ContainerImageRepo {
	return &containerImageRepo{
		client: getDockerClient("/var/run/docker.sock"),
	}
}

// List 列出镜像
func (r *containerImageRepo) List() ([]types.ContainerImage, error) {
	var resp []image.Image
	_, err := r.client.R().SetResult(&resp).SetQueryParam("all", "true").Get("/images/json")
	if err != nil {
		return nil, err
	}

	var images []types.ContainerImage
	for _, item := range resp {
		images = append(images, types.ContainerImage{
			ID:          item.ID,
			Containers:  item.Containers,
			RepoTags:    item.RepoTags,
			RepoDigests: item.RepoDigests,
			Size:        tools.FormatBytes(float64(item.Size)),
			Labels:      types.MapToKV(item.Labels),
			CreatedAt:   time.Unix(item.Created, 0),
		})
	}

	slices.SortFunc(images, func(a types.ContainerImage, b types.ContainerImage) int {
		return strings.Compare(a.ID, b.ID)
	})

	return images, nil
}

// Pull 拉取镜像
func (r *containerImageRepo) Pull(req *request.ContainerImagePull) error {
	var sb strings.Builder

	if req.Auth {
		sb.WriteString(fmt.Sprintf("docker login -u %s -p %s", req.Username, req.Password))
		if _, err := shell.Exec(sb.String()); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
		sb.Reset()
	}

	sb.WriteString(fmt.Sprintf("docker pull %s", req.Name))

	if _, err := shell.Exec(sb.String()); err != nil {
		return err
	}

	return nil
}

// Remove 删除镜像
func (r *containerImageRepo) Remove(id string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker rmi %s", id)
	return err
}

// Prune 清理未使用的镜像
func (r *containerImageRepo) Prune() error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker image prune -f")
	return err
}
