package data

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/str"
	"github.com/TheTNB/panel/pkg/types"
	"github.com/TheTNB/panel/pkg/types/docker/image"
)

type containerImageRepo struct {
	client *resty.Client
}

func NewContainerImageRepo(sock ...string) biz.ContainerImageRepo {
	if len(sock) == 0 {
		sock = append(sock, "/var/run/docker.sock")
	}
	client := resty.New()
	client.SetTimeout(1 * time.Minute)
	client.SetRetryCount(2)
	client.SetTransport(&http.Transport{
		DialContext: func(ctx context.Context, _ string, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", sock[0])
		},
	})
	client.SetBaseURL("http://d/v1.40")

	return &containerImageRepo{
		client: client,
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
			Size:        str.FormatBytes(float64(item.Size)),
			Labels:      types.MapToKV(item.Labels),
			CreatedAt:   time.Unix(item.Created, 0),
		})
	}

	return images, nil
}

// Pull 拉取镜像
func (r *containerImageRepo) Pull(req *request.ContainerImagePull) error {
	var sb strings.Builder

	if req.Auth {
		sb.WriteString(fmt.Sprintf("docker login -u %s -p %s", req.Username, req.Password))
		if _, err := shell.ExecfWithTimeout(1*time.Minute, sb.String()); err != nil { // nolint: govet
			return fmt.Errorf("login failed: %w", err)
		}
		sb.Reset()
	}

	sb.WriteString(fmt.Sprintf("docker pull %s", req.Name))

	if _, err := shell.Execf(sb.String()); err != nil { // nolint: govet
		return fmt.Errorf("pull failed: %w", err)
	}

	return nil
}

// Remove 删除镜像
func (r *containerImageRepo) Remove(id string) error {
	_, err := shell.ExecfWithTimeout(30*time.Second, "docker rmi %s", id)
	return err
}

// Prune 清理未使用的镜像
func (r *containerImageRepo) Prune() error {
	_, err := shell.ExecfWithTimeout(30*time.Second, "docker image prune -f")
	return err
}
