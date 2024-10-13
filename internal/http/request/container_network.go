package request

import "github.com/TheTNB/panel/pkg/types"

type ContainerNetworkID struct {
	ID string `json:"id" form:"id" validate:"required"`
}

type ContainerNetworkCreate struct {
	Name    string                 `form:"name" json:"name" validate:"required"`
	Driver  string                 `form:"driver" json:"driver"`
	Ipv4    types.ContainerNetwork `form:"ipv4" json:"ipv4"`
	Ipv6    types.ContainerNetwork `form:"ipv6" json:"ipv6"`
	Labels  []types.KV             `form:"labels" json:"labels"`
	Options []types.KV             `form:"options" json:"options"`
}

type ContainerNetworkConnect struct {
	Network   string `form:"network" json:"network" validate:"required"`
	Container string `form:"container" json:"container" validate:"required"`
}
