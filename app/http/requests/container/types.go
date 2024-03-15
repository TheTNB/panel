package requests

type ContainerPort struct {
	ContainerStart int    `json:"start"`
	ContainerEnd   int    `json:"end"`
	Host           string `json:"host"`
	HostStart      int    `json:"hostStart"`
	HostEnd        int    `json:"hostEnd"`
	Protocol       string `json:"protocol"`
}

type ContainerVolume struct {
	Host      string `json:"host"`
	Container string `json:"container"`
	Mode      string `json:"mode"`
}
