package data

import (
	"errors"

	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
)

type cacheRepo struct{}

func NewCacheRepo() biz.CacheRepo {
	return do.MustInvoke[biz.CacheRepo](injector)
}

func (r *cacheRepo) Get(key biz.CacheKey, defaultValue ...string) (string, error) {
	cache := new(biz.Cache)
	if err := app.Orm.Where("key = ?", key).First(cache).Error; err != nil {
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
	if err := app.Orm.Where(biz.Cache{Key: key}).FirstOrInit(cache).Error; err != nil {
		return err
	}

	cache.Value = value
	return app.Orm.Save(cache).Error
}
