package biz

import (
	"time"

	"github.com/go-rat/utils/crypt"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/http/request"
)

type DatabaseType string

const (
	DatabaseTypeMysql      DatabaseType = "mysql"
	DatabaseTypePostgresql DatabaseType = "postgresql"
	DatabaseTypeRedis      DatabaseType = "redis"
)

type DatabaseServer struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"not null;unique" json:"name"`
	Type      DatabaseType `gorm:"not null" json:"type"`
	Host      string       `gorm:"not null" json:"host"`
	Port      uint         `gorm:"not null" json:"port"`
	Username  string       `gorm:"not null" json:"username"`
	Password  string       `gorm:"not null" json:"password"`
	Remark    string       `gorm:"not null" json:"remark"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	Databases []*Database `gorm:"foreignKey:ServerID" json:"-"`
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
	Create(req *request.DatabaseServerCreate) error
	Update(req *request.DatabaseServerUpdate) error
	Delete(id uint) error
}
