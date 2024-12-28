package request

import "github.com/TheTNB/panel/pkg/types"

type ContainerNetworkID struct {
	ID string `json:"id" form:"id" validate:"required"`
}

type ContainerNetworkCreate struct {
	Name    string                          `form:"name" json:"name" validate:"required"`
	Driver  string                          `form:"driver" json:"driver" validate:"required|in:bridge,host,overlay,macvlan,ipvlan,none"`
	Ipv4    types.ContainerContainerNetwork `form:"ipv4" json:"ipv4"`
	Ipv6    types.ContainerContainerNetwork `form:"ipv6" json:"ipv6"`
	Labels  []types.KV                      `form:"labels" json:"labels"`
	Options []types.KV                      `form:"options" json:"options"`
}
