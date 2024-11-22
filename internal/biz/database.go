package biz

import (
	"time"

	"github.com/go-rat/utils/crypt"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/http/request"
)

type DatabaseStatus string

const (
	DatabaseStatusNormal  DatabaseStatus = "normal"
	DatabaseStatusInvalid DatabaseStatus = "invalid"
)

type Database struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ServerID  uint           `gorm:"not null" json:"server_id"`
	Name      string         `gorm:"not null" json:"name"`
	Status    DatabaseStatus `gorm:"not null" json:"status"`
	Username  string         `gorm:"not null" json:"username"`
	Password  string         `gorm:"not null" json:"password"`
	Remark    string         `gorm:"not null" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	Server *DatabaseServer `gorm:"foreignKey:ServerID" json:"server"`
}

func (r *Database) BeforeSave(tx *gorm.DB) error {
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

func (r *Database) AfterFind(tx *gorm.DB) error {
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

type DatabaseRepo interface {
	Count() (int64, error)
	List(page, limit uint) ([]*Database, int64, error)
	Get(id uint) (*Database, error)
	Create(req *request.DatabaseCreate) error
	Update(req *request.DatabaseUpdate) error
	Delete(id uint) error
	Add(serverID uint, name string) error
}
