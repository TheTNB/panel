package biz

import (
	"time"

	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/ssh"
)

type SSH struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	Name      string           `gorm:"not null" json:"name"`
	Host      string           `gorm:"not null" json:"host"`
	Port      uint             `gorm:"not null" json:"port"`
	Config    ssh.ClientConfig `gorm:"not null;serializer:json" json:"config"`
	Remark    string           `gorm:"not null" json:"remark"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type SSHRepo interface {
	List(page, limit uint) ([]*SSH, int64, error)
	Get(id uint) (*SSH, error)
	Create(req *request.SSHCreate) error
	Update(req *request.SSHUpdate) error
	Delete(id uint) error
}
