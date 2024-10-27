package data

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
)

type containerImageRepo struct {
	client *client.Client
}

func NewContainerImageRepo(sock ...string) biz.ContainerImageRepo {
	if len(sock) == 0 {
		sock = append(sock, "/run/podman/podman.sock")
	}
	cli, _ := client.NewClientWithOpts(client.WithHost("unix://"+sock[0]), client.WithAPIVersionNegotiation())
	return &containerImageRepo{
		client: cli,
	}
}

// List 列出镜像
func (r *containerImageRepo) List() ([]image.Summary, error) {
	return r.client.ImageList(context.Background(), image.ListOptions{
		All: true,
	})
}

// Pull 拉取镜像
func (r *containerImageRepo) Pull(req *request.ContainerImagePull) error {
	options := image.PullOptions{}
	if req.Auth {
		authConfig := registry.AuthConfig{
			Username: req.Username,
			Password: req.Password,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		options.RegistryAuth = authStr
	}

	out, err := r.client.ImagePull(context.Background(), req.Name, options)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(io.Discard, out)
	return err
}

// Remove 删除镜像
func (r *containerImageRepo) Remove(id string) error {
	_, err := r.client.ImageRemove(context.Background(), id, image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	return err
}

// Prune 清理未使用的镜像
func (r *containerImageRepo) Prune() error {
	_, err := r.client.ImagesPrune(context.Background(), filters.NewArgs())
	return err
}
