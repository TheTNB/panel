package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"

	"github.com/TheTNB/panel/v2/pkg/types"
)

type ContainerCreate struct {
	Name            string                  `form:"name" json:"name"`
	Image           string                  `form:"image" json:"image"`
	Ports           []types.ContainerPort   `form:"ports" json:"ports"`
	Network         string                  `form:"network" json:"network"`
	Volumes         []types.ContainerVolume `form:"volumes" json:"volumes"`
	Labels          []types.KV              `form:"labels" json:"labels"`
	Env             []types.KV              `form:"env" json:"env"`
	Entrypoint      []string                `form:"entrypoint" json:"entrypoint"`
	Command         []string                `form:"command" json:"command"`
	RestartPolicy   string                  `form:"restart_policy" json:"restart_policy"`
	AutoRemove      bool                    `form:"auto_remove" json:"auto_remove"`
	Privileged      bool                    `form:"privileged" json:"privileged"`
	OpenStdin       bool                    `form:"openStdin" json:"open_stdin"`
	PublishAllPorts bool                    `form:"publish_all_ports" json:"publish_all_ports"`
	Tty             bool                    `form:"tty" json:"tty"`
	CPUShares       int64                   `form:"cpu_shares" json:"cpu_shares"`
	CPUs            int64                   `form:"cpus" json:"cpus"`
	Memory          int64                   `form:"memory" json:"memory"`
}

func (r *ContainerCreate) Authorize(ctx http.Context) error {
	return nil
}

func (r *ContainerCreate) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":  "required|string",
		"image": "required|string",
		"ports": "slice",
		/*"ports.*.host":            "string",
		"ports.*.host_start":      "int",
		"ports.*.host_end":        "int",
		"ports.*.container_start": "int",
		"ports.*.container_end":   "int",
		"ports.*.protocol":        "string|in:tcp,udp",*/
		"network": "string",
		"volumes": "slice",
		/*"volumes.*.host":          "string",
		"volumes.*.container":     "string",
		"volumes.*.mode":          "string|in:ro,rw",*/
		"labels":            "slice",
		"env":               "slice",
		"entrypoint":        "slice",
		"command":           "slice",
		"restart_policy":    "string|in:always,on-failure,unless-stopped,no",
		"auto_remove":       "bool",
		"privileged":        "bool",
		"open_stdin":        "bool",
		"publish_all_ports": "bool",
		"tty":               "bool",
		"cpu_shares":        "int",
		"cpus":              "int",
		"memory":            "int",
	}
}

func (r *ContainerCreate) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"cpu_shares": "int",
		"cpus":       "int",
		"memory":     "int",
	}
}

func (r *ContainerCreate) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ContainerCreate) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ContainerCreate) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
