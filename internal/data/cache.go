package data

import (
	"encoding/json"
	"errors"
	"slices"

	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/pkg/api"
	"github.com/TheTNB/panel/pkg/apploader"
)

type cacheRepo struct {
	api *api.API
}

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

func (r *cacheRepo) UpdateApps() error {
	remote, err := r.api.Apps()
	if err != nil {
		return err
	}

	// 去除本地不存在的应用
	*remote = slices.Clip(slices.DeleteFunc(*remote, func(app *api.App) bool {
		_, err = apploader.Get(app.Slug)
		return err != nil
	}))

	encoded, err := json.Marshal(remote)
	if err != nil {
		return err
	}

	return r.Set(biz.CacheKeyApps, string(encoded))
}

func (r *cacheRepo) UpdateRewrites() error {
	rewrites, err := r.api.RewritesByType("nginx")
	if err != nil {
		return err
	}

	encoded, err := json.Marshal(rewrites)
	if err != nil {
		return err
	}

	return r.Set(biz.CacheKeyRewrites, string(encoded))
}
