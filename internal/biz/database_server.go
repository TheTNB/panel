package biz

import (
	"github.com/TheTNB/panel/internal/app"
	"github.com/go-rat/utils/crypt"
	"time"

	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/http/request"
)

type DatabaseServerStatus string

const (
	DatabaseServerStatusValid   DatabaseServerStatus = "valid"
	DatabaseServerStatusInvalid DatabaseServerStatus = "invalid"
)

type DatabaseServer struct {
	ID        uint                 `gorm:"primaryKey" json:"id"`
	Name      string               `gorm:"not null;unique" json:"name"`
	Type      DatabaseType         `gorm:"not null" json:"type"`
	Host      string               `gorm:"not null" json:"host"`
	Port      uint                 `gorm:"not null" json:"port"`
	Username  string               `gorm:"not null" json:"username"`
	Password  string               `gorm:"not null" json:"password"`
	Status    DatabaseServerStatus `gorm:"-:all" json:"status"`
	Remark    string               `gorm:"not null" json:"remark"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

func (r *DatabaseServer) BeforeSave(tx *gorm.DB) error {
	crypter, err := crypt.NewXChacha20Poly1305([]byte(app.Key))
	if err != nil {
		return err
	}

	r.Password, err = crypter.Encrypt([]byte(r.Password))
	if err != nil {
		return err
	}

	return nil

}

func (r *DatabaseServer) AfterFind(tx *gorm.DB) error {
	crypter, err := crypt.NewXChacha20Poly1305([]byte(app.Key))
	if err != nil {
		return err
	}

	password, err := crypter.Decrypt(r.Password)
	if err == nil {
		r.Password = string(password)
	}

	return nil
}

type DatabaseServerRepo interface {
	Count() (int64, error)
	List(page, limit uint) ([]*DatabaseServer, int64, error)
	Get(id uint) (*DatabaseServer, error)
	GetByName(name string) (*DatabaseServer, error)
	Create(req *request.DatabaseServerCreate) error
	Update(req *request.DatabaseServerUpdate) error
	UpdateRemark(req *request.DatabaseServerUpdateRemark) error
	Delete(id uint) error
	ClearUsers(id uint) error
	Sync(id uint) error
}
