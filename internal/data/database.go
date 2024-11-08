package data

import (
	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
)

type databaseRepo struct{}

func NewDatabaseRepo() biz.DatabaseRepo {
	return &databaseRepo{}
}

func (d databaseRepo) Count() (int64, error) {
	var count int64
	if err := app.Orm.Model(&biz.Database{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (d databaseRepo) List(page, limit uint) ([]*biz.Database, int64, error) {
	var database []*biz.Database
	var total int64
	err := app.Orm.Model(&biz.Database{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&database).Error
	return database, total, err
}

func (d databaseRepo) Get(id uint) (*biz.Database, error) {
	database := new(biz.Database)
	if err := app.Orm.Where("id = ?", id).First(database).Error; err != nil {
		return nil, err
	}

	return database, nil
}

func (d databaseRepo) Create(req *request.DatabaseCreate) error {
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

func (d databaseRepo) Update(req *request.DatabaseUpdate) error {
	database := &biz.Database{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
		Remark:   req.Remark,
	}

	return app.Orm.Model(database).Where("id = ?", req.ID).Updates(database).Error
}

func (d databaseRepo) Delete(id uint) error {
	return app.Orm.Delete(&biz.Database{}, id).Error
}
