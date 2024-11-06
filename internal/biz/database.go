package biz

import (
	"errors"
	"time"

	"github.com/go-rat/utils/crypt"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
)

type Database struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;unique" json:"name"`
	Type      string    `gorm:"not null" json:"type"`
	Host      string    `gorm:"not null" json:"host"`
	Port      int       `gorm:"not null" json:"port"`
	Username  string    `gorm:"not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Remark    string    `gorm:"not null" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DatabaseItems []*DatabaseItem `json:"-"`
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

func (r *Database) BeforeDelete(tx *gorm.DB) error {
	if r.Name == "local_mysql" && !app.IsCli {
		return errors.New("can't delete local_mysql, if you must delete it, please uninstall mysql")
	}
	if r.Name == "local_postgresql" && !app.IsCli {
		return errors.New("can't delete local_postgresql, if you must delete it, please uninstall postgresql")
	}
	if r.Name == "local_redis" && !app.IsCli {
		return errors.New("can't delete local_redis, if you must delete it, please uninstall redis")
	}

	return nil
}
