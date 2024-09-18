package data

import (
	"errors"

	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/internal/panel"
)

type monitorRepo struct {
	settingRepo biz.SettingRepo
}

func NewMonitorRepo() biz.MonitorRepo {
	return &monitorRepo{
		settingRepo: NewSettingRepo(),
	}
}

func (r monitorRepo) GetSetting() (*request.MonitorSetting, error) {
	monitor, err := r.settingRepo.Get(biz.SettingKeyMonitor)
	if err != nil {
		return nil, err
	}
	monitorDays, err := r.settingRepo.Get(biz.SettingKeyMonitorDays)
	if err != nil {
		return nil, err
	}

	setting := new(request.MonitorSetting)
	setting.Enabled = cast.ToBool(monitor)
	setting.Days = cast.ToInt(monitorDays)

	return setting, nil
}

func (r monitorRepo) UpdateSetting(setting *request.MonitorSetting) error {
	if err := r.settingRepo.Set(biz.SettingKeyMonitor, cast.ToString(setting.Enabled)); err != nil {
		return err
	}
	if err := r.settingRepo.Set(biz.SettingKeyMonitorDays, cast.ToString(setting.Days)); err != nil {
		return err
	}

	return nil
}

func (r monitorRepo) Clear() error {
	return panel.Orm.Delete(&biz.Monitor{}).Error
}

func (r monitorRepo) List(start, end carbon.Carbon) ([]*biz.Monitor, error) {
	var monitors []*biz.Monitor
	if err := panel.Orm.Where("created_at BETWEEN ? AND ?", start, end).Find(&monitors).Error; err != nil {
		return nil, err
	}

	if len(monitors) == 0 {
		return nil, errors.New("没有找到数据")
	}

	return monitors, nil
}
