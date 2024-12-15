package biz

import (
	"time"

	"gorm.io/gorm"

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

func (r *SSH) BeforeSave(tx *gorm.DB) error {
	// TODO fix
	/*var err error
	r.Config.Key, err = app.Crypter.Encrypt([]byte(r.Config.Key))
	if err != nil {
		return err
	}
	r.Config.Password, err = app.Crypter.Encrypt([]byte(r.Config.Password))
	if err != nil {
		return err
	}*/

	return nil

}

func (r *SSH) AfterFind(tx *gorm.DB) error {
	// TODO fix
	/*key, err := app.Crypter.Decrypt(r.Config.Key)
	if err == nil {
		r.Config.Key = string(key)
	}
	password, err := app.Crypter.Decrypt(r.Config.Password)
	if err == nil {
		r.Config.Password = string(password)
	}*/

	return nil
}

type SSHRepo interface {
	List(page, limit uint) ([]*SSH, int64, error)
	Get(id uint) (*SSH, error)
	Create(req *request.SSHCreate) error
	Update(req *request.SSHUpdate) error
	Delete(id uint) error
}
