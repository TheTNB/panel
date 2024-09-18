package data

import (
	"errors"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/panel"
)

type cacheRepo struct{}

func NewCacheRepo() biz.CacheRepo {
	return &cacheRepo{}
}

func (r *cacheRepo) Get(key biz.CacheKey, defaultValue ...string) (string, error) {
	cache := new(biz.Cache)
	if err := panel.Orm.Where("key = ?", key).First(cache).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	}

	if cache.Value == "" && len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return cache.Value, nil
}

func (r *cacheRepo) Set(key biz.CacheKey, value string) error {
	cache := new(biz.Cache)
	if err := panel.Orm.Where(biz.Cache{Key: key}).FirstOrInit(cache).Error; err != nil {
		return err
	}

	cache.Value = value
	return panel.Orm.Save(cache).Error
}
