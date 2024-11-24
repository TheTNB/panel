package biz

import (
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/types"
)

type DatabaseRepo interface {
	Count() (int64, error)
	List(page, limit uint) ([]types.Database, int64, error)
	Create(req *request.DatabaseCreate) error
	Delete(serverID uint, name string) error
}
