package biz

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"not null;unique" json:"username"`
	Password  string         `gorm:"not null" json:"password"`
	Email     string         `gorm:"not null" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type UserRepo interface {
	CheckPassword(username, password string) (*User, error)
	Get(id uint) (*User, error)
	Save(user *User) error
}
