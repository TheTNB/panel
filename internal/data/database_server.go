package data

import (
	"fmt"

	"github.com/samber/do/v2"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/db"
)

type databaseServerRepo struct{}

func NewDatabaseServerRepo() biz.DatabaseServerRepo {
	return do.MustInvoke[biz.DatabaseServerRepo](injector)
}

func (r databaseServerRepo) Count() (int64, error) {
	var count int64
	if err := app.Orm.Model(&biz.DatabaseServer{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r databaseServerRepo) List(page, limit uint) ([]*biz.DatabaseServer, int64, error) {
	var databaseServer []*biz.DatabaseServer
	var total int64
	err := app.Orm.Model(&biz.DatabaseServer{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&databaseServer).Error
	return databaseServer, total, err
}

func (r databaseServerRepo) Get(id uint) (*biz.DatabaseServer, error) {
	databaseServer := new(biz.DatabaseServer)
	if err := app.Orm.Where("id = ?", id).First(databaseServer).Error; err != nil {
		return nil, err
	}

	return databaseServer, nil
}

func (r databaseServerRepo) Create(req *request.DatabaseServerCreate) error {
	switch biz.DatabaseType(req.Type) {
	case biz.DatabaseTypeMysql:
		if _, err := db.NewMySQL(req.Username, req.Password, fmt.Sprintf("%s:%d", req.Host, req.Port)); err != nil {
			return err
		}
	case biz.DatabaseTypePostgresql:
		if _, err := db.NewPostgres(req.Username, req.Password, req.Host, req.Port); err != nil {
			return err
		}
	case biz.DatabaseTypeRedis:
		if _, err := db.NewRedis(req.Username, req.Password, fmt.Sprintf("%s:%d", req.Host, req.Port)); err != nil {
			return err
		}

	}

	databaseServer := &biz.DatabaseServer{
		Name:     req.Name,
		Type:     biz.DatabaseType(req.Type),
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Remark:   req.Remark,
	}

	return app.Orm.Create(databaseServer).Error
}

func (r databaseServerRepo) Update(req *request.DatabaseServerUpdate) error {
	switch biz.DatabaseType(req.Type) {
	case biz.DatabaseTypeMysql:
		if _, err := db.NewMySQL(req.Username, req.Password, fmt.Sprintf("%s:%d", req.Host, req.Port)); err != nil {
			return err
		}
	case biz.DatabaseTypePostgresql:
		if _, err := db.NewPostgres(req.Username, req.Password, req.Host, req.Port); err != nil {
			return err
		}
	case biz.DatabaseTypeRedis:
		if _, err := db.NewRedis(req.Username, req.Password, fmt.Sprintf("%s:%d", req.Host, req.Port)); err != nil {
			return err
		}

	}

	return app.Orm.Model(&biz.DatabaseServer{}).Where("id = ?", req.ID).Select("*").Updates(&biz.DatabaseServer{
		Name:     req.Name,
		Type:     biz.DatabaseType(req.Type),
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Remark:   req.Remark,
	}).Error
}

func (r databaseServerRepo) Delete(id uint) error {
	ds := new(biz.DatabaseServer)
	if err := app.Orm.Where("id = ?", id).First(ds).Error; err != nil {
		return err
	}

	return app.Orm.Delete(&biz.DatabaseServer{}, id).Error
}
