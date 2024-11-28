package biz

import "time"

type CacheKey string

const (
	CacheKeyApps     CacheKey = "apps"
	CacheKeyRewrites CacheKey = "rewrites"
)

type Cache struct {
	Key       CacheKey  `gorm:"primaryKey" json:"key"`
	Value     string    `gorm:"not null" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CacheRepo interface {
	Get(key CacheKey, defaultValue ...string) (string, error)
	Set(key CacheKey, value string) error
	UpdateApps() error
	UpdateRewrites() error
}
