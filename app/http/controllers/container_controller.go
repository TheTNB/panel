package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/goravel/framework/contracts/http"
	requests "panel/app/http/requests/container"
	"panel/internal/services"
)

type ContainerController struct {
	container services.Container
}

func NewContainerController() *ContainerController {
	return &ContainerController{
		container: services.NewContainer(),
	}
}

func (r *ContainerController) ContainerList(ctx http.Context) http.Response {
	containers, err := r.container.ContainerListAll()
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, containers)
}

func (r *ContainerController) ContainerSearch(ctx http.Context) http.Response {
	fields := strings.Fields(ctx.Request().Query("names"))
	containers, err := r.container.ContainerListByNames(fields)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, containers)
}

func (r *ContainerController) ContainerCreate(ctx http.Context) http.Response {
	var request requests.ContainerCreate
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	var hostConf container.HostConfig
	var networkConf network.NetworkingConfig

	portMap := make(nat.PortMap)
	for _, port := range request.Ports {
		if port.ContainerStart-port.ContainerEnd != port.HostStart-port.HostEnd {
			return Error(ctx, http.StatusUnprocessableEntity, fmt.Sprintf("容器端口和主机端口数量不匹配（容器: %d 主机: %d）", port.ContainerStart-port.ContainerEnd, port.HostStart-port.HostEnd))
		}
		if port.ContainerStart > port.ContainerEnd || port.HostStart > port.HostEnd || port.ContainerStart < 1 || port.HostStart < 1 {
			return Error(ctx, http.StatusUnprocessableEntity, "端口范围不正确")
		}

		count := 0
		for host := port.HostStart; host <= port.HostEnd; host++ {
			bindItem := nat.PortBinding{HostPort: strconv.Itoa(host), HostIP: port.Host}
			portMap[nat.Port(fmt.Sprintf("%d/%s", port.ContainerStart+count, port.Protocol))] = []nat.PortBinding{bindItem}
			count++
		}
	}

	exposed := make(nat.PortSet)
	for port := range portMap {
		exposed[port] = struct{}{}
	}

	if request.Network != "" {
		switch request.Network {
		case "host", "none", "bridge":
			hostConf.NetworkMode = container.NetworkMode(request.Network)
		}
		networkConf.EndpointsConfig = map[string]*network.EndpointSettings{request.Network: {}}
	} else {
		networkConf = network.NetworkingConfig{}
	}

	hostConf.Privileged = request.Privileged
	hostConf.AutoRemove = request.AutoRemove
	hostConf.CPUShares = request.CPUShares
	hostConf.PublishAllPorts = request.PublishAllPorts
	hostConf.RestartPolicy = container.RestartPolicy{Name: container.RestartPolicyMode(request.RestartPolicy)}
	if request.RestartPolicy == "on-failure" {
		hostConf.RestartPolicy.MaximumRetryCount = 5
	}
	hostConf.NanoCPUs = request.CPUs * 1000000000
	hostConf.Memory = request.Memory * 1024 * 1024
	hostConf.MemorySwap = 0
	hostConf.PortBindings = portMap
	hostConf.Binds = []string{}

	volumes := make(map[string]struct{})
	for _, volume := range request.Volumes {
		volumes[volume.Container] = struct{}{}
		hostConf.Binds = append(hostConf.Binds, fmt.Sprintf("%s:%s:%s", volume.Host, volume.Container, volume.Mode))
	}

	id, err := r.container.ContainerCreate(request.Name,
		container.Config{
			Image:        request.Image,
			Env:          request.Env,
			Entrypoint:   request.Entrypoint,
			Cmd:          request.Command,
			Labels:       r.container.SliceToMap(request.Labels),
			ExposedPorts: exposed,
			OpenStdin:    request.OpenStdin,
			Tty:          request.Tty,
			Volumes:      volumes,
		},
		hostConf,
		networkConf,
	)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, id)
}

func (r *ContainerController) ContainerRemove(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerRemove(request.ID); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

func (r *ContainerController) ContainerStart(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerStart(request.ID); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

func (r *ContainerController) ContainerStop(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerStop(request.ID); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

func (r *ContainerController) ContainerRestart(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerRestart(request.ID); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

func (r *ContainerController) ContainerPause(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerPause(request.ID); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

func (r *ContainerController) ContainerUnpause(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerUnpause(request.ID); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

func (r *ContainerController) ContainerInspect(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.ContainerInspect(request.ID)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, data)
}

func (r *ContainerController) ContainerKill(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerKill(request.ID); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

func (r *ContainerController) ContainerRename(ctx http.Context) http.Response {
	var request requests.ContainerRename
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerRename(request.ID, request.Name); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

func (r *ContainerController) ContainerStats(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.ContainerStats(request.ID)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, data)
}

func (r *ContainerController) ContainerExist(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	exist, err := r.container.ContainerExist(request.ID)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, exist)
}

func (r *ContainerController) ContainerLogs(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := Sanitize(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.ContainerLogs(request.ID)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, data)
}

func (r *ContainerController) ContainerPrune(ctx http.Context) http.Response {
	if err := r.container.ContainerPrune(); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

func (r *ContainerController) NetworkList(ctx http.Context) http.Response {
	networks, err := r.container.NetworkList()
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, networks)
}

func (r *ContainerController) NetworkCreate(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) NetworkRemove(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) NetworkExist(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) NetworkInspect(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) NetworkConnect(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) NetworkDisconnect(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) NetworkPrune(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) ImageList(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) ImageExist(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) ImagePull(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) ImageRemove(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) ImagePrune(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) ImageInspect(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) VolumeList(ctx http.Context) http.Response {
	volumes, err := r.container.VolumeList()
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, volumes)
}

func (r *ContainerController) VolumeCreate(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) VolumeExist(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) VolumeInspect(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) VolumeRemove(ctx http.Context) http.Response {
	return Success(ctx, nil)
}

func (r *ContainerController) VolumePrune(ctx http.Context) http.Response {
	return Success(ctx, nil)
}
