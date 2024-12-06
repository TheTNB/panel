package data

import (
	"fmt"
	"log/slog"
	"slices"

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

	for server := range slices.Values(databaseServer) {
		r.checkServer(server)
	}

	return databaseServer, total, err
}

func (r databaseServerRepo) Get(id uint) (*biz.DatabaseServer, error) {
	databaseServer := new(biz.DatabaseServer)
	if err := app.Orm.Where("id = ?", id).First(databaseServer).Error; err != nil {
		return nil, err
	}

	r.checkServer(databaseServer)

	return databaseServer, nil
}

func (r databaseServerRepo) GetByName(name string) (*biz.DatabaseServer, error) {
	databaseServer := new(biz.DatabaseServer)
	if err := app.Orm.Where("name = ?", name).First(databaseServer).Error; err != nil {
		return nil, err
	}

	r.checkServer(databaseServer)

	return databaseServer, nil
}

func (r databaseServerRepo) Create(req *request.DatabaseServerCreate) error {
	databaseServer := &biz.DatabaseServer{
		Name:     req.Name,
		Type:     biz.DatabaseType(req.Type),
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Remark:   req.Remark,
	}

	if !r.checkServer(databaseServer) {
		return fmt.Errorf("check server connection failed")
	}

	return app.Orm.Create(databaseServer).Error
}

func (r databaseServerRepo) Update(req *request.DatabaseServerUpdate) error {
	server, err := r.Get(req.ID)
	if err != nil {
		return err
	}

	server.Name = req.Name
	server.Host = req.Host
	server.Port = req.Port
	server.Username = req.Username
	server.Password = req.Password
	server.Remark = req.Remark

	if !r.checkServer(server) {
		return fmt.Errorf("check server connection failed")
	}

	return app.Orm.Save(server).Error
}

func (r databaseServerRepo) UpdateRemark(req *request.DatabaseServerUpdateRemark) error {
	return app.Orm.Model(&biz.DatabaseServer{}).Where("id = ?", req.ID).Update("remark", req.Remark).Error
}

func (r databaseServerRepo) Delete(id uint) error {
	// 删除服务器下的所有用户
	if err := NewDatabaseUserRepo().DeleteByServerID(id); err != nil {
		return err
	}

	return app.Orm.Where("id = ?", id).Delete(&biz.DatabaseServer{}).Error
}

func (r databaseServerRepo) Sync(id uint) error {
	server, err := r.Get(id)
	if err != nil {
		return err
	}

	users := make([]*biz.DatabaseUser, 0)
	if err = app.Orm.Where("server_id = ?", id).Find(&users).Error; err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		defer mysql.Close()
		allUsers, err := mysql.Users()
		if err != nil {
			return err
		}
		for user := range slices.Values(allUsers) {
			if !slices.ContainsFunc(users, func(a *biz.DatabaseUser) bool {
				return a.Username == user.User && a.Host == user.Host
			}) && !slices.Contains([]string{"root", "mysql.sys", "mysql.session", "mysql.infoschema"}, user.User) {
				newUser := &biz.DatabaseUser{
					ServerID: id,
					Username: user.User,
					Host:     user.Host,
					Remark:   fmt.Sprintf("sync from server %s", server.Name),
				}
				if err = app.Orm.Create(newUser).Error; err != nil {
					app.Logger.Warn("sync database user failed", slog.Any("err", err))
				}
			}
		}
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		defer postgres.Close()
		allUsers, err := postgres.Users()
		if err != nil {
			return err
		}
		for user := range slices.Values(allUsers) {
			if !slices.ContainsFunc(users, func(a *biz.DatabaseUser) bool {
				return a.Username == user.Role
			}) && !slices.Contains([]string{"postgres"}, user.Role) {
				newUser := &biz.DatabaseUser{
					ServerID: id,
					Username: user.Role,
					Remark:   fmt.Sprintf("sync from server %s", server.Name),
				}
				app.Orm.Create(newUser)
			}
		}
	}

	return nil
}

// checkServer 检查服务器连接
func (r databaseServerRepo) checkServer(server *biz.DatabaseServer) bool {
	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err == nil {
			_ = mysql.Close()
			server.Status = biz.DatabaseServerStatusValid
			return true
		}
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err == nil {
			_ = postgres.Close()
			server.Status = biz.DatabaseServerStatusValid
			return true
		}
	case biz.DatabaseTypeRedis:
		redis, err := db.NewRedis(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err == nil {
			_ = redis.Close()
			server.Status = biz.DatabaseServerStatusValid
			return true
		}
	}

	server.Status = biz.DatabaseServerStatusInvalid
	return false
}
