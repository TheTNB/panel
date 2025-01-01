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
	"github.com/tnb-labs/panel/pkg/types/docker/container"
)

type containerRepo struct {
	client *resty.Client
}

func NewContainerRepo() biz.ContainerRepo {
	return &containerRepo{
		client: getDockerClient("/var/run/docker.sock"),
	}
}

// ListAll 列出所有容器
func (r *containerRepo) ListAll() ([]types.Container, error) {
	var resp []container.Container
	_, err := r.client.R().SetResult(&resp).SetQueryParam("all", "true").Get("/containers/json")
	if err != nil {
		return nil, err
	}

	var containers []types.Container
	for _, item := range resp {
		ports := make([]types.ContainerPort, 0)
		for _, port := range item.Ports {
			ports = append(ports, types.ContainerPort{
				ContainerStart: uint(port.PrivatePort),
				ContainerEnd:   uint(port.PublicPort),
				HostStart:      uint(port.PublicPort),
				HostEnd:        uint(port.PublicPort),
				Protocol:       port.Type,
				Host:           port.IP,
			})
		}
		if len(item.Names) == 0 {
			item.Names = append(item.Names, "")
		}
		containers = append(containers, types.Container{
			ID:        item.ID,
			Name:      strings.TrimPrefix(item.Names[0], "/"), // https://github.com/moby/moby/issues/7519
			Image:     item.Image,
			ImageID:   item.ImageID,
			Command:   item.Command,
			CreatedAt: time.Unix(item.Created, 0),
			State:     item.State,
			Status:    item.Status,
			Ports:     ports,
			Labels:    types.MapToKV(item.Labels),
		})
	}

	slices.SortFunc(containers, func(a types.Container, b types.Container) int {
		return strings.Compare(a.Name, b.Name)
	})

	return containers, nil
}

// ListByName 根据名称搜索容器
func (r *containerRepo) ListByName(names string) ([]types.Container, error) {
	containers, err := r.ListAll()
	if err != nil {
		return nil, err
	}

	containers = slices.DeleteFunc(containers, func(item types.Container) bool {
		return !strings.Contains(item.Name, names)
	})

	return containers, nil
}

// Create 创建容器
func (r *containerRepo) Create(req *request.ContainerCreate) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("docker run -d --name %s", req.Name))
	if req.PublishAllPorts {
		sb.WriteString(" -P")
	} else {
		for _, port := range req.Ports {
			sb.WriteString(" -p ")
			if port.Host != "" {
				sb.WriteString(fmt.Sprintf("%s:", port.Host))
			}
			if port.HostStart == port.HostEnd || port.ContainerStart == port.ContainerEnd {
				sb.WriteString(fmt.Sprintf("%d:%d/%s", port.HostStart, port.ContainerStart, port.Protocol))
			} else {
				sb.WriteString(fmt.Sprintf("%d-%d:%d-%d/%s", port.HostStart, port.HostEnd, port.ContainerStart, port.ContainerEnd, port.Protocol))
			}
		}
	}
	if req.Network != "" {
		sb.WriteString(fmt.Sprintf(" --network %s", req.Network))
	}
	for _, volume := range req.Volumes {
		sb.WriteString(fmt.Sprintf(" -v %s:%s:%s", volume.Host, volume.Container, volume.Mode))
	}
	for _, label := range req.Labels {
		sb.WriteString(fmt.Sprintf(" --label %s=%s", label.Key, label.Value))
	}
	for _, env := range req.Env {
		sb.WriteString(fmt.Sprintf(" -e %s=%s", env.Key, env.Value))
	}
	if len(req.Entrypoint) > 0 {
		sb.WriteString(fmt.Sprintf(" --entrypoint '%s'", strings.Join(req.Entrypoint, " ")))
	}
	if len(req.Command) > 0 {
		sb.WriteString(fmt.Sprintf(" '%s'", strings.Join(req.Command, " ")))
	}
	if req.RestartPolicy != "" {
		sb.WriteString(fmt.Sprintf(" --restart %s", req.RestartPolicy))
	}
	if req.AutoRemove {
		sb.WriteString(" --rm")
	}
	if req.Privileged {
		sb.WriteString(" --privileged")
	}
	if req.OpenStdin {
		sb.WriteString(" -i")
	}
	if req.Tty {
		sb.WriteString(" -t")
	}
	if req.CPUShares > 0 {
		sb.WriteString(fmt.Sprintf(" --cpu-shares %d", req.CPUShares))
	}
	if req.CPUs > 0 {
		sb.WriteString(fmt.Sprintf(" --cpus %d", req.CPUs))
	}
	if req.Memory > 0 {
		sb.WriteString(fmt.Sprintf(" --memory %d", req.Memory))
	}

	sb.WriteString(" %s bash")
	return shell.ExecfWithTTY(sb.String(), req.Image)
}

// Remove 移除容器
func (r *containerRepo) Remove(id string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker rm -f %s", id)
	return err
}

// Start 启动容器
func (r *containerRepo) Start(id string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker start %s", id)
	return err
}

// Stop 停止容器
func (r *containerRepo) Stop(id string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker stop %s", id)
	return err
}

// Restart 重启容器
func (r *containerRepo) Restart(id string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker restart %s", id)
	return err
}

// Pause 暂停容器
func (r *containerRepo) Pause(id string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker pause %s", id)
	return err
}

// Unpause 恢复容器
func (r *containerRepo) Unpause(id string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker unpause %s", id)
	return err
}

// Kill 杀死容器
func (r *containerRepo) Kill(id string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker kill %s", id)
	return err
}

// Rename 重命名容器
func (r *containerRepo) Rename(id string, newName string) error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker rename %s %s", id, newName)
	return err
}

// Logs 查看容器日志
func (r *containerRepo) Logs(id string) (string, error) {
	return shell.ExecfWithTimeout(120*time.Second, "docker logs %s", id)
}

// Prune 清理未使用的容器
func (r *containerRepo) Prune() error {
	_, err := shell.ExecfWithTimeout(120*time.Second, "docker container prune -f")
	return err
}
