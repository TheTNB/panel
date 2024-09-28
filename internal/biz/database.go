package biz

import "time"

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
}
