package requests

type ContainerPort struct {
	ContainerStart int    `form:"start" json:"start"`
	ContainerEnd   int    `form:"end" json:"end"`
	Host           string `form:"host" json:"host"`
	HostStart      int    `form:"host_start" json:"host_start"`
	HostEnd        int    `form:"host_end" json:"host_end"`
	Protocol       string `form:"protocol" json:"protocol"`
}

type ContainerVolume struct {
	Host      string `form:"start" json:"host"`
	Container string `form:"start" json:"container"`
	Mode      string `form:"start" json:"mode"`
}

type Network struct {
	Enabled bool   `form:"enabled" json:"enabled"`
	Gateway string `form:"gateway" json:"gateway"`
	IPRange string `form:"ip_range" json:"ip_range"`
	Subnet  string `form:"subnet" json:"subnet"`
}
