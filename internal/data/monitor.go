package data

import (
	"errors"
	"time"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/tnb-labs/panel/internal/biz"
	"github.com/tnb-labs/panel/internal/http/request"
)

type monitorRepo struct {
	db      *gorm.DB
	setting biz.SettingRepo
}

func NewMonitorRepo(db *gorm.DB, setting biz.SettingRepo) biz.MonitorRepo {
	return &monitorRepo{
		db:      db,
		setting: setting,
	}
}

func (r monitorRepo) GetSetting() (*request.MonitorSetting, error) {
	monitor, err := r.setting.Get(biz.SettingKeyMonitor)
	if err != nil {
		return nil, err
	}
	monitorDays, err := r.setting.Get(biz.SettingKeyMonitorDays)
	if err != nil {
		return nil, err
	}

	setting := new(request.MonitorSetting)
	setting.Enabled = cast.ToBool(monitor)
	setting.Days = cast.ToUint(monitorDays)

	return setting, nil
}

func (r monitorRepo) UpdateSetting(setting *request.MonitorSetting) error {
	if err := r.setting.Set(biz.SettingKeyMonitor, cast.ToString(setting.Enabled)); err != nil {
		return err
	}
	if err := r.setting.Set(biz.SettingKeyMonitorDays, cast.ToString(setting.Days)); err != nil {
		return err
	}

	return nil
}

func (r monitorRepo) Clear() error {
	return r.db.Where("1 = 1").Delete(&biz.Monitor{}).Error
}

func (r monitorRepo) List(start, end time.Time) ([]*biz.Monitor, error) {
	var monitors []*biz.Monitor
	if err := r.db.Where("created_at BETWEEN ? AND ?", start, end).Find(&monitors).Error; err != nil {
		return nil, err
	}

	if len(monitors) == 0 {
		return nil, errors.New("没有找到数据")
	}

	return monitors, nil
}
