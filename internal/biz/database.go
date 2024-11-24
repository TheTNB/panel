package biz

import (
	"github.com/TheTNB/panel/internal/http/request"
)

type DatabaseStatus string

const (
	DatabaseStatusValid   DatabaseStatus = "valid"
	DatabaseStatusInvalid DatabaseStatus = "invalid"
)

type Database struct {
	Name     string         `json:"name"`
	ServerID uint           `json:"server_id"`
	Status   DatabaseStatus `json:"status"`
	Remark   string         `json:"remark"`
}

type DatabaseRepo interface {
	List(page, limit uint) ([]*Database, int64, error)
	Create(req *request.DatabaseCreate) error
	Delete(serverID uint, name string) error
}
