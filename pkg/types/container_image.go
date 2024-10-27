package types

import (
	"time"
)

type ContainerImage struct {
	ID         string    `json:"id"`
	Containers int64     `json:"containers"`
	Tag        string    `json:"tag"`
	Size       string    `json:"size"`
	CreatedAt  time.Time `json:"created_at"`
}
