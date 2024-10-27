package types

import "time"

type ContainerPort struct {
	IP          string `json:"ip,omitempty"`
	PrivatePort uint   `json:"private_port"`
	PublicPort  uint   `json:"public_port,omitempty"`
	Type        string `json:"type"`
}

type Container struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Image     string          `json:"image"`
	ImageID   string          `json:"image_id"`
	Command   string          `json:"command"`
	CreatedAt time.Time       `json:"created_at"`
	Ports     []ContainerPort `json:"ports"`
	Labels    []KV
	State     string
	Status    string
}

type ContainerCreatePort struct {
	ContainerStart int    `form:"container_start" json:"container_start"`
	ContainerEnd   int    `form:"container_end" json:"container_end"`
	Host           string `form:"host" json:"host"`
	HostStart      int    `form:"host_start" json:"host_start"`
	HostEnd        int    `form:"host_end" json:"host_end"`
	Protocol       string `form:"protocol" json:"protocol"`
}

type ContainerVolume struct {
	Host      string `form:"host" json:"host"`
	Container string `form:"container" json:"container"`
	Mode      string `form:"mode" json:"mode"`
}

type ContainerNetwork struct {
	Enabled bool   `form:"enabled" json:"enabled"`
	Gateway string `form:"gateway" json:"gateway"`
	IPRange string `form:"ip_range" json:"ip_range"`
	Subnet  string `form:"subnet" json:"subnet"`
}
