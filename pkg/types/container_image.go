package types

import (
	"time"
)

type ContainerImage struct {
	ID          string    `json:"id"`
	Containers  int64     `json:"containers"`
	RepoTags    []string  `json:"repo_tags"`
	RepoDigests []string  `json:"repo_digests"`
	Size        string    `json:"size"`
	Labels      []KV      `json:"labels"`
	CreatedAt   time.Time `json:"created_at"`
}
