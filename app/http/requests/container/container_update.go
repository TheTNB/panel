package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"

	"github.com/TheTNB/panel/v2/pkg/types"
)

type ContainerUpdate struct {
	ID              string                  `form:"id" json:"id"`
	Name            string                  `form:"name" json:"name"`
	Image           string                  `form:"image" json:"image"`
	Ports           []types.ContainerPort   `form:"ports" json:"ports"`
	Network         string                  `form:"network" json:"network"`
	Volumes         []types.ContainerVolume `form:"volumes" json:"volumes"`
	Labels          []string                `form:"labels" json:"labels"`
	Env             []string                `form:"env" json:"env"`
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

func (r *ContainerUpdate) Authorize(ctx http.Context) error {
	return nil
}

func (r *ContainerUpdate) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id":                "required|string",
		"name":              "required|string",
		"image":             "required|string",
		"ports":             "slice",
		"network":           "string",
		"volumes":           "slice",
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

func (r *ContainerUpdate) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ContainerUpdate) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ContainerUpdate) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ContainerUpdate) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
