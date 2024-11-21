package data

import (
	"fmt"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/db"
)

type databaseRepo struct {
	databaseServer biz.DatabaseServerRepo
}

func NewDatabaseRepo() biz.DatabaseRepo {
	return &databaseRepo{
		databaseServer: NewDatabaseServerRepo(),
	}
}

func (r databaseRepo) Count() (int64, error) {
	var count int64
	if err := app.Orm.Model(&biz.Database{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r databaseRepo) List(page, limit uint) ([]*biz.Database, int64, error) {
	var database []*biz.Database
	var total int64
	err := app.Orm.Model(&biz.Database{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&database).Error
	return database, total, err
}

func (r databaseRepo) Get(id uint) (*biz.Database, error) {
	database := new(biz.Database)
	if err := app.Orm.Where("id = ?", id).First(database).Error; err != nil {
		return nil, err
	}

	return database, nil
}

func (r databaseRepo) Create(req *request.DatabaseCreate) error {
	server, err := r.databaseServer.Get(req.ServerID)
	if err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		if err = mysql.UserCreate(req.Username, req.Password); err != nil {
			return err
		}
		if err = mysql.DatabaseCreate(req.Name); err != nil {
			return err
		}
		if err = mysql.PrivilegesGrant(req.Username, req.Name); err != nil {
			return err
		}
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		if err = postgres.UserCreate(req.Username, req.Password); err != nil {
			return err
		}
		if err = postgres.DatabaseCreate(req.Name); err != nil {
			return err
		}
		if err = postgres.PrivilegesGrant(req.Username, req.Name); err != nil {
			return err
		}
	}

	database := &biz.Database{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
		ServerID: req.ServerID,
		Status:   biz.DatabaseStatusInvalid,
		Remark:   req.Remark,
	}

	return app.Orm.Create(database).Error
}

func (r databaseRepo) Update(req *request.DatabaseUpdate) error {
	database := &biz.Database{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
		Remark:   req.Remark,
	}

	return app.Orm.Model(database).Where("id = ?", req.ID).Omit("ServerID").Updates(database).Error
}

func (r databaseRepo) Delete(id uint) error {
	return app.Orm.Delete(&biz.Database{}, id).Error
}
