package biz

import (
	"time"

	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/ssh"
)

type SSH struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	Host      string           `json:"host"`
	Port      uint             `json:"port"`
	Config    ssh.ClientConfig `json:"config"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type SSHRepo interface {
	GetInfo() (map[string]any, error)
	UpdateInfo(req *request.SSHUpdateInfo) error
}
