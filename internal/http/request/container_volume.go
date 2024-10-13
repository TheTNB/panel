package request

import "github.com/TheTNB/panel/pkg/types"

type ContainerVolumeID struct {
	ID string `json:"id" form:"id" validate:"required"`
}

type ContainerVolumeCreate struct {
	Name    string     `form:"name" json:"name" validate:"required"`
	Driver  string     `form:"driver" json:"driver"`
	Labels  []types.KV `form:"labels" json:"labels"`
	Options []types.KV `form:"options" json:"options"`
}
