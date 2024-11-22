package data

import (
	"errors"
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
)

type monitorRepo struct{}

func NewMonitorRepo() biz.MonitorRepo {
	return do.MustInvoke[biz.MonitorRepo](injector)
}

func (r monitorRepo) GetSetting() (*request.MonitorSetting, error) {
	repo := NewSettingRepo()
	monitor, err := repo.Get(biz.SettingKeyMonitor)
	if err != nil {
		return nil, err
	}
	monitorDays, err := repo.Get(biz.SettingKeyMonitorDays)
	if err != nil {
		return nil, err
	}

	setting := new(request.MonitorSetting)
	setting.Enabled = cast.ToBool(monitor)
	setting.Days = cast.ToInt(monitorDays)

	return setting, nil
}

func (r monitorRepo) UpdateSetting(setting *request.MonitorSetting) error {
	repo := NewSettingRepo()
	if err := repo.Set(biz.SettingKeyMonitor, cast.ToString(setting.Enabled)); err != nil {
		return err
	}
	if err := repo.Set(biz.SettingKeyMonitorDays, cast.ToString(setting.Days)); err != nil {
		return err
	}

	return nil
}

func (r monitorRepo) Clear() error {
	return app.Orm.Where("1 = 1").Delete(&biz.Monitor{}).Error
}

func (r monitorRepo) List(start, end time.Time) ([]*biz.Monitor, error) {
	var monitors []*biz.Monitor
	if err := app.Orm.Where("created_at BETWEEN ? AND ?", start, end).Find(&monitors).Error; err != nil {
		return nil, err
	}

	if len(monitors) == 0 {
		return nil, errors.New("没有找到数据")
	}

	return monitors, nil
}
