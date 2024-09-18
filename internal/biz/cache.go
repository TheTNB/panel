package biz

import "github.com/golang-module/carbon/v2"

type CacheKey string

const (
	CacheKeyApps     CacheKey = "apps"
	CacheKeyRewrites CacheKey = "rewrites"
)

type Cache struct {
	Key       CacheKey        `gorm:"primaryKey" json:"key"`
	Value     string          `gorm:"not null" json:"value"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type CacheRepo interface {
	Get(key CacheKey, defaultValue ...string) (string, error)
	Set(key CacheKey, value string) error
}
