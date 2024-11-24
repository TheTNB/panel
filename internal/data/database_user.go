package data

import (
	"fmt"
	"slices"

	"github.com/samber/do/v2"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/db"
)

type databaseUserRepo struct{}

func NewDatabaseUserRepo() biz.DatabaseUserRepo {
	return do.MustInvoke[biz.DatabaseUserRepo](injector)
}

func (r databaseUserRepo) Count() (int64, error) {
	var count int64
	if err := app.Orm.Model(&biz.DatabaseUser{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r databaseUserRepo) List(page, limit uint) ([]*biz.DatabaseUser, int64, error) {
	var user []*biz.DatabaseUser
	var total int64
	err := app.Orm.Model(&biz.DatabaseUser{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&user).Error

	for u := range slices.Values(user) {
		r.fillUser(u)
	}

	return user, total, err
}

func (r databaseUserRepo) Get(id uint) (*biz.DatabaseUser, error) {
	user := new(biz.DatabaseUser)
	if err := app.Orm.Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}

	r.fillUser(user)

	return user, nil
}

func (r databaseUserRepo) Create(req *request.DatabaseUserCreate) error {
	server, err := NewDatabaseServerRepo().Get(req.ServerID)
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
		for name := range slices.Values(req.Privileges) {
			if err = mysql.PrivilegesGrant(req.Username, name); err != nil {
				return err
			}
		}
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		if err = postgres.UserCreate(req.Username, req.Password); err != nil {
			return err
		}
		for name := range slices.Values(req.Privileges) {
			if err = postgres.PrivilegesGrant(req.Username, name); err != nil {
				return err
			}
		}
	}

	user := &biz.DatabaseUser{
		ServerID: req.ServerID,
		Username: req.Username,
		Password: req.Password,
		Remark:   req.Remark,
	}

	return app.Orm.Create(user).Error
}

func (r databaseUserRepo) Update(req *request.DatabaseUserUpdate) error {
	user, err := r.Get(req.ID)
	if err != nil {
		return err
	}

	server, err := NewDatabaseServerRepo().Get(user.ServerID)
	if err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		if req.Password != "" {
			if err = mysql.UserPassword(user.Username, req.Password); err != nil {
				return err
			}
		}
		for name := range slices.Values(req.Privileges) {
			if err = mysql.PrivilegesGrant(user.Username, name); err != nil {
				return err
			}
		}
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		if req.Password != "" {
			if err = postgres.UserPassword(user.Username, req.Password); err != nil {
				return err
			}
		}
		for name := range slices.Values(req.Privileges) {
			if err = postgres.PrivilegesGrant(user.Username, name); err != nil {
				return err
			}
		}
	}

	user.Password = req.Password
	user.Remark = req.Remark

	return app.Orm.Save(user).Error
}

func (r databaseUserRepo) Delete(id uint) error {
	user, err := r.Get(id)
	if err != nil {
		return err
	}

	server, err := NewDatabaseServerRepo().Get(user.ServerID)
	if err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		_ = mysql.UserDrop(user.Username)
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		_ = postgres.DatabaseDrop(user.Username)
	}

	return app.Orm.Where("id = ?", id).Delete(&biz.DatabaseUser{}).Error
}

func (r databaseUserRepo) DeleteByServerID(serverID uint) error {
	return app.Orm.Where("server_id = ?", serverID).Delete(&biz.DatabaseUser{}).Error
}

func (r databaseUserRepo) fillUser(user *biz.DatabaseUser) {
	server, err := NewDatabaseServerRepo().Get(user.ServerID)
	if err == nil {
		switch server.Type {
		case biz.DatabaseTypeMysql:
			mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
			if err == nil {
				privileges, _ := mysql.UserPrivileges(user.Username, user.Host)
				user.Privileges = privileges
			}
			if _, err := db.NewMySQL(user.Username, user.Password, fmt.Sprintf("%s:%d", server.Host, server.Port)); err == nil {
				user.Status = biz.DatabaseUserStatusValid
			} else {
				user.Status = biz.DatabaseUserStatusInvalid
			}
		case biz.DatabaseTypePostgresql:
			postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
			if err == nil {
				privileges, _ := postgres.UserPrivileges(user.Username)
				user.Privileges = privileges
			}
			if _, err := db.NewPostgres(user.Username, user.Password, server.Host, server.Port); err == nil {
				user.Status = biz.DatabaseUserStatusValid
			} else {
				user.Status = biz.DatabaseUserStatusInvalid
			}
		}
	}
	// 初始化，防止 nil
	if user.Privileges == nil {
		user.Privileges = make(map[string][]string)
	}
}
