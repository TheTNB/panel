package biz

import (
	"time"

	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/http/request"
)

type DatabaseUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ServerID  uint      `gorm:"not null" json:"server_id"`
	Username  string    `gorm:"not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Remark    string    `gorm:"not null" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *DatabaseUser) BeforeSave(tx *gorm.DB) error {
	var err error
	r.Password, err = app.Crypter.Encrypt([]byte(r.Password))
	if err != nil {
		return err
	}

	return nil

}

func (r *DatabaseUser) AfterFind(tx *gorm.DB) error {
	password, err := app.Crypter.Decrypt(r.Password)
	if err == nil {
		r.Password = string(password)
	}

	return nil
}

type DatabaseUserRepo interface {
	Count() (int64, error)
	List(page, limit uint) ([]*DatabaseUser, int64, error)
	Get(id uint) (*DatabaseUser, error)
	Create(req *request.DatabaseUserCreate) error
	Update(req *request.DatabaseUserUpdate) error
	Delete(id uint) error
	Sync(id uint) error
}
