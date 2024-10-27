package data

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/types"
)

type containerRepo struct {
	cmd string
}

func NewContainerRepo(cmd ...string) biz.ContainerRepo {
	if len(cmd) == 0 {
		cmd = append(cmd, "docker")
	}
	return &containerRepo{
		cmd: cmd[0],
	}
}

// ListAll 列出所有容器
func (r *containerRepo) ListAll() ([]types.Container, error) {
	output, err := shell.ExecfWithTimeout(10*time.Second, "%s ps -a --format json", r.cmd)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(output, "\n")

	var containers []types.Container
	for _, line := range lines {
		if line == "" {
			continue
		}

		var item struct {
			Command      string `json:"Command"`
			CreatedAt    string `json:"CreatedAt"`
			ID           string `json:"ID"`
			Image        string `json:"Image"`
			Labels       string `json:"Labels"`
			LocalVolumes string `json:"LocalVolumes"`
			Mounts       string `json:"Mounts"`
			Names        string `json:"Names"`
			Networks     string `json:"Networks"`
			Ports        string `json:"Ports"`
			RunningFor   string `json:"RunningFor"`
			Size         string `json:"Size"`
			State        string `json:"State"`
			Status       string `json:"Status"`
		}
		if err = json.Unmarshal([]byte(line), &item); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}

		createdAt, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", item.CreatedAt)
		containers = append(containers, types.Container{
			ID:        item.ID,
			Name:      item.Names,
			Image:     item.Image,
			Command:   item.Command,
			CreatedAt: createdAt,
			Ports:     r.parsePorts(item.Ports),
			Labels:    types.SliceToKV(strings.Split(item.Labels, ",")),
			State:     item.State,
			Status:    item.Status,
		})
	}

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
	sb.WriteString(fmt.Sprintf("%s create --name %s --image %s", r.cmd, req.Name, req.Image))

	for _, port := range req.Ports {
		sb.WriteString(fmt.Sprintf(" -p %s:%d-%d:%d-%d/%s", port.Host, port.HostStart, port.HostEnd, port.ContainerStart, port.ContainerEnd, port.Protocol))
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
	if req.PublishAllPorts {
		sb.WriteString(" -P")
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

	return shell.ExecfWithTimeout(10*time.Second, sb.String()) // nolint: govet
}

// Remove 移除容器
func (r *containerRepo) Remove(id string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s rm -f %s", r.cmd, id)
	return err
}

// Start 启动容器
func (r *containerRepo) Start(id string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s start %s", r.cmd, id)
	return err
}

// Stop 停止容器
func (r *containerRepo) Stop(id string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s stop %s", r.cmd, id)
	return err
}

// Restart 重启容器
func (r *containerRepo) Restart(id string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s restart %s", r.cmd, id)
	return err
}

// Pause 暂停容器
func (r *containerRepo) Pause(id string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s pause %s", r.cmd, id)
	return err
}

// Unpause 恢复容器
func (r *containerRepo) Unpause(id string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s unpause %s", r.cmd, id)
	return err
}

// Kill 杀死容器
func (r *containerRepo) Kill(id string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s kill %s", r.cmd, id)
	return err
}

// Rename 重命名容器
func (r *containerRepo) Rename(id string, newName string) error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s rename %s %s", r.cmd, id, newName)
	return err
}

// Logs 查看容器日志
func (r *containerRepo) Logs(id string) (string, error) {
	return shell.ExecfWithTimeout(10*time.Second, "%s logs %s", r.cmd, id)
}

// Prune 清理未使用的容器
func (r *containerRepo) Prune() error {
	_, err := shell.ExecfWithTimeout(10*time.Second, "%s container prune -f", r.cmd)
	return err
}

func (r *containerRepo) parsePorts(ports string) []types.ContainerPort {
	var portList []types.ContainerPort

	re := regexp.MustCompile(`(?P<host>[\d.:]+)?:(?P<public>\d+)->(?P<private>\d+)/(?P<protocol>\w+)`)

	entries := strings.Split(ports, ", ") // 0.0.0.0:3306->3306/tcp, :::3306->3306/tcp, 33060/tcp
	for _, entry := range entries {
		matches := re.FindStringSubmatch(entry)
		if len(matches) == 0 {
			continue
		}

		host := matches[1]
		public := matches[2]
		private := matches[3]
		protocol := matches[4]

		portList = append(portList, types.ContainerPort{
			Host:           host,
			HostStart:      cast.ToUint(public),
			HostEnd:        cast.ToUint(public),
			ContainerStart: cast.ToUint(private),
			ContainerEnd:   cast.ToUint(private),
			Protocol:       protocol,
		})
	}

	return portList
}
