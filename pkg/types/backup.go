package types

import "time"

type BackupFile struct {
	Name string    `json:"name"`
	Path string    `json:"path"`
	Size string    `json:"size"`
	Time time.Time `json:"time"`
}
