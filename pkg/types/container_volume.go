package types

import (
	"time"
)

type ContainerVolume struct {
	Name       string    `json:"name"`
	Driver     string    `json:"driver"`
	Scope      string    `json:"scope"`
	MountPoint string    `json:"mount_point"`
	CreatedAt  time.Time `json:"created_at"`
	Labels     []KV      `json:"labels"`
	Options    []KV      `json:"options"`
	RefCount   int64     `json:"ref_count"`
	Size       string    `json:"size"`
}
