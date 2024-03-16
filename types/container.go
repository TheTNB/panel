package types

type ContainerPort struct {
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
