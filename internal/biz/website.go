package biz

import (
	"time"

	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/types"
)

type Website struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;unique" json:"name"`
	Status    bool      `gorm:"not null;default:true" json:"status"`
	Path      string    `gorm:"not null" json:"path"`
	PHP       int       `gorm:"not null" json:"php"`
	SSL       bool      `gorm:"not null" json:"ssl"`
	Remark    string    `gorm:"not null" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Cert *Cert `gorm:"foreignKey:WebsiteID" json:"cert"`
}

type WebsiteRepo interface {
	UpdateDefaultConfig(req *request.WebsiteDefaultConfig) error
	Count() (int64, error)
	Get(id uint) (*types.WebsiteSetting, error)
	GetByName(name string) (*types.WebsiteSetting, error)
	List(page, limit uint) ([]*Website, int64, error)
	Create(req *request.WebsiteCreate) (*Website, error)
	Update(req *request.WebsiteUpdate) error
	Delete(req *request.WebsiteDelete) error
	ClearLog(id uint) error
	UpdateRemark(id uint, remark string) error
	ResetConfig(id uint) error
	UpdateStatus(id uint, status bool) error
}
