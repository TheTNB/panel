package biz

import "github.com/TheTNB/panel/internal/http/request"

type SSHRepo interface {
	GetInfo() (map[string]any, error)
	UpdateInfo(req *request.SSHUpdateInfo) error
}
