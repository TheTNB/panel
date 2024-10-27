package types

import (
	"time"
)

type Container struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Image     string          `json:"image"`
	ImageID   string          `json:"image_id"`
	Command   string          `json:"command"`
	State     string          `json:"state"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	Ports     []ContainerPort `json:"ports"`
	Labels    []KV            `json:"labels"`
}

type ContainerPort struct {
	ContainerStart uint   `form:"container_start" json:"container_start"`
	ContainerEnd   uint   `form:"container_end" json:"container_end"`
	Host           string `form:"host" json:"host"`
	HostStart      uint   `form:"host_start" json:"host_start"`
	HostEnd        uint   `form:"host_end" json:"host_end"`
	Protocol       string `form:"protocol" json:"protocol"`
}

type ContainerContainerVolume struct {
	Host      string `form:"host" json:"host"`
	Container string `form:"container" json:"container"`
	Mode      string `form:"mode" json:"mode"`
}

type ContainerContainerNetwork struct {
	Enabled bool   `form:"enabled" json:"enabled"`
	Gateway string `form:"gateway" json:"gateway"`
	IPRange string `form:"ip_range" json:"ip_range"`
	Subnet  string `form:"subnet" json:"subnet"`
}
