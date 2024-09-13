package data

import (
	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
)

type cronRepo struct{}

func NewCronRepo() biz.CronRepo {
	return &cronRepo{}
}

func (r *cronRepo) Count() (int64, error) {
	var count int64
	if err := app.Orm.Model(&biz.Cron{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
