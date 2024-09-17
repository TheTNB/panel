package request

import "github.com/TheTNB/panel/pkg/types"

type ContainerVolumeCreate struct {
	Name    string     `form:"name" json:"name"`
	Driver  string     `form:"driver" json:"driver"`
	Labels  []types.KV `form:"labels" json:"labels"`
	Options []types.KV `form:"options" json:"options"`
}
