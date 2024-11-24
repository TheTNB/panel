package biz

import (
	"github.com/TheTNB/panel/internal/http/request"
)

type DatabaseStatus string

type Database struct {
	Name     string `json:"name"`
	Server   string `json:"server"`
	Encoding string `json:"encoding"`
}

type DatabaseRepo interface {
	List(page, limit uint) ([]*Database, int64, error)
	Create(req *request.DatabaseCreate) error
	Delete(serverID uint, name string) error
}
