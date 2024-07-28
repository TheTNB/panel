package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"

	requests "github.com/TheTNB/panel/v2/app/http/requests/container"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/str"
)

type ContainerController struct {
	container services.Container
}

func NewContainerController() *ContainerController {
	return &ContainerController{
		container: services.NewContainer(),
	}
}

// ContainerList
//
//	@Summary		获取容器列表
//	@Description	获取所有容器列表
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/list [get]
func (r *ContainerController) ContainerList(ctx http.Context) http.Response {
	containers, err := r.container.ContainerListAll()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	paged, total := h.Paginate(ctx, containers)

	items := make([]any, 0)
	for _, item := range paged {
		var name string
		if len(item.Names) > 0 {
			name = item.Names[0]
		}
		items = append(items, map[string]any{
			"id":       item.ID,
			"name":     strings.TrimLeft(name, "/"),
			"image":    item.Image,
			"image_id": item.ImageID,
			"command":  item.Command,
			"created":  carbon.FromTimestamp(item.Created).ToDateTimeString(),
			"ports":    item.Ports,
			"labels":   item.Labels,
			"state":    item.State,
			"status":   item.Status,
		})
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": items,
	})
}

// ContainerSearch
//
//	@Summary		搜索容器
//	@Description	根据容器名称搜索容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			name	query		string	true	"容器名称"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/search [get]
func (r *ContainerController) ContainerSearch(ctx http.Context) http.Response {
	fields := strings.Fields(ctx.Request().Query("name"))
	containers, err := r.container.ContainerListByNames(fields)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, containers)
}

// ContainerCreate
//
//	@Summary		创建容器
//	@Description	创建一个容器
//	@Tags			容器
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ContainerCreate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/create [post]
func (r *ContainerController) ContainerCreate(ctx http.Context) http.Response {
	var request requests.ContainerCreate
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	var hostConf container.HostConfig
	var networkConf network.NetworkingConfig

	portMap := make(nat.PortMap)
	for _, port := range request.Ports {
		if port.ContainerStart-port.ContainerEnd != port.HostStart-port.HostEnd {
			return h.Error(ctx, http.StatusUnprocessableEntity, fmt.Sprintf("容器端口和主机端口数量不匹配（容器: %d 主机: %d）", port.ContainerStart-port.ContainerEnd, port.HostStart-port.HostEnd))
		}
		if port.ContainerStart > port.ContainerEnd || port.HostStart > port.HostEnd || port.ContainerStart < 1 || port.HostStart < 1 {
			return h.Error(ctx, http.StatusUnprocessableEntity, "端口范围不正确")
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
	for _, v := range request.Volumes {
		volumes[v.Container] = struct{}{}
		hostConf.Binds = append(hostConf.Binds, fmt.Sprintf("%s:%s:%s", v.Host, v.Container, v.Mode))
	}

	id, err := r.container.ContainerCreate(request.Name,
		container.Config{
			Image:        request.Image,
			Env:          r.container.KVToSlice(request.Env),
			Entrypoint:   request.Entrypoint,
			Cmd:          request.Command,
			Labels:       r.container.KVToMap(request.Labels),
			ExposedPorts: exposed,
			OpenStdin:    request.OpenStdin,
			Tty:          request.Tty,
			Volumes:      volumes,
		},
		hostConf,
		networkConf,
	)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = r.container.ContainerStart(id); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, id)
}

// ContainerRemove
//
//	@Summary		删除容器
//	@Description	删除一个容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/remove [post]
func (r *ContainerController) ContainerRemove(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerRemove(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ContainerStart
//
//	@Summary		启动容器
//	@Description	启动一个容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/start [post]
func (r *ContainerController) ContainerStart(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerStart(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ContainerStop
//
//	@Summary		停止容器
//	@Description	停止一个容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/stop [post]
func (r *ContainerController) ContainerStop(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerStop(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ContainerRestart
//
//	@Summary		重启容器
//	@Description	重启一个容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/restart [post]
func (r *ContainerController) ContainerRestart(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerRestart(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ContainerPause
//
//	@Summary		暂停容器
//	@Description	暂停一个容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
func (r *ContainerController) ContainerPause(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerPause(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ContainerUnpause
//
//	@Summary		取消暂停容器
//	@Description	取消暂停一个容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//
//	@Router			/panel/container/unpause [post]
func (r *ContainerController) ContainerUnpause(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerUnpause(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ContainerInspect
//
//	@Summary		查看容器
//	@Description	查看一个容器的详细信息
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/inspect [get]
func (r *ContainerController) ContainerInspect(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.ContainerInspect(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, data)
}

// ContainerKill
//
//	@Summary		杀死容器
//	@Description	杀死一个容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/kill [post]
func (r *ContainerController) ContainerKill(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerKill(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ContainerRename
//
//	@Summary		重命名容器
//	@Description	重命名一个容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ContainerRename	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/rename [post]
func (r *ContainerController) ContainerRename(ctx http.Context) http.Response {
	var request requests.ContainerRename
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ContainerRename(request.ID, request.Name); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ContainerStats
//
//	@Summary		查看容器状态
//	@Description	查看一个容器的状态信息
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/stats [get]
func (r *ContainerController) ContainerStats(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.ContainerStats(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, data)
}

// ContainerExist
//
//	@Summary		检查容器是否存在
//	@Description	检查一个容器是否存在
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/exist [get]
func (r *ContainerController) ContainerExist(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	exist, err := r.container.ContainerExist(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, exist)
}

// ContainerLogs
//
//	@Summary		查看容器日志
//	@Description	查看一个容器的日志
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/logs [get]
func (r *ContainerController) ContainerLogs(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.ContainerLogs(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, data)
}

// ContainerPrune
//
//	@Summary		清理容器
//	@Description	清理无用的容器
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/container/prune [post]
func (r *ContainerController) ContainerPrune(ctx http.Context) http.Response {
	if err := r.container.ContainerPrune(); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// NetworkList
//
//	@Summary		获取网络列表
//	@Description	获取所有网络列表
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/network/list [get]
func (r *ContainerController) NetworkList(ctx http.Context) http.Response {
	networks, err := r.container.NetworkList()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	paged, total := h.Paginate(ctx, networks)

	items := make([]any, 0)
	for _, item := range paged {
		var ipamConfig []any
		for _, v := range item.IPAM.Config {
			ipamConfig = append(ipamConfig, map[string]any{
				"subnet":      v.Subnet,
				"gateway":     v.Gateway,
				"ip_range":    v.IPRange,
				"aux_address": v.AuxAddress,
			})
		}
		items = append(items, map[string]any{
			"id":         item.ID,
			"name":       item.Name,
			"driver":     item.Driver,
			"ipv6":       item.EnableIPv6,
			"scope":      item.Scope,
			"internal":   item.Internal,
			"attachable": item.Attachable,
			"ingress":    item.Ingress,
			"labels":     item.Labels,
			"options":    item.Options,
			"ipam": map[string]any{
				"config":  ipamConfig,
				"driver":  item.IPAM.Driver,
				"options": item.IPAM.Options,
			},
			"created": carbon.FromStdTime(item.Created).ToDateTimeString(),
		})
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": items,
	})
}

// NetworkCreate
//
//	@Summary		创建网络
//	@Description	创建一个网络
//	@Tags			容器
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.NetworkCreate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/network/create [post]
func (r *ContainerController) NetworkCreate(ctx http.Context) http.Response {
	var request requests.NetworkCreate
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	id, err := r.container.NetworkCreate(request)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, id)
}

// NetworkRemove
//
//	@Summary		删除网络
//	@Description	删除一个网络
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/network/remove [post]
func (r *ContainerController) NetworkRemove(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.NetworkRemove(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// NetworkExist
//
//	@Summary		检查网络是否存在
//	@Description	检查一个网络是否存在
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/network/exist [get]
func (r *ContainerController) NetworkExist(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	exist, err := r.container.NetworkExist(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, exist)
}

// NetworkInspect
//
//	@Summary		查看网络
//	@Description	查看一个网络的详细信息
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/network/inspect [get]
func (r *ContainerController) NetworkInspect(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.NetworkInspect(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, data)
}

// NetworkConnect
//
//	@Summary		连接容器到网络
//	@Description	连接一个容器到一个网络
//	@Tags			容器
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.NetworkConnectDisConnect	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/network/connect [post]
func (r *ContainerController) NetworkConnect(ctx http.Context) http.Response {
	var request requests.NetworkConnectDisConnect
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.NetworkConnect(request.Network, request.Container); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// NetworkDisconnect
//
//	@Summary		从网络断开容器
//	@Description	从一个网络断开一个容器
//	@Tags			容器
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.NetworkConnectDisConnect	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/network/disconnect [post]
func (r *ContainerController) NetworkDisconnect(ctx http.Context) http.Response {
	var request requests.NetworkConnectDisConnect
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.NetworkDisconnect(request.Network, request.Container); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// NetworkPrune
//
//	@Summary		清理网络
//	@Description	清理无用的网络
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/container/network/prune [post]
func (r *ContainerController) NetworkPrune(ctx http.Context) http.Response {
	if err := r.container.NetworkPrune(); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ImageList
//
//	@Summary		获取镜像列表
//	@Description	获取所有镜像列表
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/image/list [get]
func (r *ContainerController) ImageList(ctx http.Context) http.Response {
	images, err := r.container.ImageList()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	paged, total := h.Paginate(ctx, images)

	items := make([]any, 0)
	for _, item := range paged {
		items = append(items, map[string]any{
			"id":           item.ID,
			"created":      carbon.FromTimestamp(item.Created).ToDateTimeString(),
			"containers":   item.Containers,
			"size":         str.FormatBytes(float64(item.Size)),
			"labels":       item.Labels,
			"repo_tags":    item.RepoTags,
			"repo_digests": item.RepoDigests,
		})
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": items,
	})
}

// ImageExist
//
//	@Summary		检查镜像是否存在
//	@Description	检查一个镜像是否存在
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/image/exist [get]
func (r *ContainerController) ImageExist(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	exist, err := r.container.ImageExist(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, exist)
}

// ImagePull
//
//	@Summary		拉取镜像
//	@Description	拉取一个镜像
//	@Tags			容器
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ImagePull	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/image/pull [post]
func (r *ContainerController) ImagePull(ctx http.Context) http.Response {
	var request requests.ImagePull
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ImagePull(request); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ImageRemove
//
//	@Summary		删除镜像
//	@Description	删除一个镜像
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/image/remove [post]
func (r *ContainerController) ImageRemove(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.ImageRemove(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ImagePrune
//
//	@Summary		清理镜像
//	@Description	清理无用的镜像
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/container/image/prune [post]
func (r *ContainerController) ImagePrune(ctx http.Context) http.Response {
	if err := r.container.ImagePrune(); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ImageInspect
//
//	@Summary		查看镜像
//	@Description	查看一个镜像的详细信息
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/image/inspect [get]
func (r *ContainerController) ImageInspect(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.ImageInspect(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, data)
}

// VolumeList
//
//	@Summary		获取卷列表
//	@Description	获取所有卷列表
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/volume/list [get]
func (r *ContainerController) VolumeList(ctx http.Context) http.Response {
	volumes, err := r.container.VolumeList()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	paged, total := h.Paginate(ctx, volumes)

	items := make([]any, 0)
	for _, item := range paged {
		var usage any
		if item.UsageData != nil {
			usage = map[string]any{
				"ref_count": item.UsageData.RefCount,
				"size":      str.FormatBytes(float64(item.UsageData.Size)),
			}
		}
		items = append(items, map[string]any{
			"id":      item.Name,
			"created": carbon.Parse(item.CreatedAt).ToDateTimeString(),
			"driver":  item.Driver,
			"mount":   item.Mountpoint,
			"labels":  item.Labels,
			"options": item.Options,
			"scope":   item.Scope,
			"status":  item.Status,
			"usage":   usage,
		})
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": items,
	})
}

// VolumeCreate
//
//	@Summary		创建卷
//	@Description	创建一个卷
//	@Tags			容器
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.VolumeCreate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/volume/create [post]
func (r *ContainerController) VolumeCreate(ctx http.Context) http.Response {
	var request requests.VolumeCreate
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.VolumeCreate(request)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, data.Name)
}

// VolumeExist
//
//	@Summary		检查卷是否存在
//	@Description	检查一个卷是否存在
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/volume/exist [get]
func (r *ContainerController) VolumeExist(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	exist, err := r.container.VolumeExist(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, exist)
}

// VolumeInspect
//
//	@Summary		查看卷
//	@Description	查看一个卷的详细信息
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/volume/inspect [get]
func (r *ContainerController) VolumeInspect(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	data, err := r.container.VolumeInspect(request.ID)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, data)
}

// VolumeRemove
//
//	@Summary		删除卷
//	@Description	删除一个卷
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.ID	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/container/volume/remove [post]
func (r *ContainerController) VolumeRemove(ctx http.Context) http.Response {
	var request requests.ID
	if sanitize := h.SanitizeRequest(ctx, &request); sanitize != nil {
		return sanitize
	}

	if err := r.container.VolumeRemove(request.ID); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// VolumePrune
//
//	@Summary		清理卷
//	@Description	清理无用的卷
//	@Tags			容器
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/container/volume/prune [post]
func (r *ContainerController) VolumePrune(ctx http.Context) http.Response {
	if err := r.container.VolumePrune(); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}
