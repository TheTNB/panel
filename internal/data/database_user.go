package data

import (
	"fmt"
	"slices"

	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/db"
)

type databaseUserRepo struct {
	db     *gorm.DB
	server biz.DatabaseServerRepo
}

func NewDatabaseUserRepo(db *gorm.DB, server biz.DatabaseServerRepo) biz.DatabaseUserRepo {
	return &databaseUserRepo{
		db:     db,
		server: server,
	}
}

func (r databaseUserRepo) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&biz.DatabaseUser{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r databaseUserRepo) List(page, limit uint) ([]*biz.DatabaseUser, int64, error) {
	var user []*biz.DatabaseUser
	var total int64
	err := r.db.Model(&biz.DatabaseUser{}).Preload("Server").Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&user).Error

	for u := range slices.Values(user) {
		r.fillUser(u)
	}

	return user, total, err
}

func (r databaseUserRepo) Get(id uint) (*biz.DatabaseUser, error) {
	user := new(biz.DatabaseUser)
	if err := r.db.Preload("Server").Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}

	r.fillUser(user)

	return user, nil
}

func (r databaseUserRepo) Create(req *request.DatabaseUserCreate) error {
	server, err := r.server.Get(req.ServerID)
	if err != nil {
		return err
	}

	user := new(biz.DatabaseUser)
	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		defer mysql.Close()
		if err = mysql.UserCreate(req.Username, req.Password, req.Host); err != nil {
			return err
		}
		for name := range slices.Values(req.Privileges) {
			if err = mysql.PrivilegesGrant(req.Username, name, req.Host); err != nil {
				return err
			}
		}
		user = &biz.DatabaseUser{
			ServerID: req.ServerID,
			Username: req.Username,
			Host:     req.Host,
		}
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		defer postgres.Close()
		if err = postgres.UserCreate(req.Username, req.Password); err != nil {
			return err
		}
		for name := range slices.Values(req.Privileges) {
			if err = postgres.PrivilegesGrant(req.Username, name); err != nil {
				return err
			}
		}
		user = &biz.DatabaseUser{
			ServerID: req.ServerID,
			Username: req.Username,
		}
	}

	if err = r.db.FirstOrInit(user, user).Error; err != nil {
		return err
	}

	user.Password = req.Password
	user.Remark = req.Remark

	return r.db.Save(user).Error
}

func (r databaseUserRepo) Update(req *request.DatabaseUserUpdate) error {
	user, err := r.Get(req.ID)
	if err != nil {
		return err
	}

	server, err := r.server.Get(user.ServerID)
	if err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		defer mysql.Close()
		if req.Password != "" {
			if err = mysql.UserPassword(user.Username, req.Password, user.Host); err != nil {
				return err
			}
		}
		for name := range slices.Values(req.Privileges) {
			if err = mysql.PrivilegesGrant(user.Username, name, user.Host); err != nil {
				return err
			}
		}
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		defer postgres.Close()
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

	return r.db.Save(user).Error
}

func (r databaseUserRepo) UpdateRemark(req *request.DatabaseUserUpdateRemark) error {
	user, err := r.Get(req.ID)
	if err != nil {
		return err
	}

	user.Remark = req.Remark

	return r.db.Save(user).Error
}

func (r databaseUserRepo) Delete(id uint) error {
	user, err := r.Get(id)
	if err != nil {
		return err
	}

	server, err := r.server.Get(user.ServerID)
	if err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		defer mysql.Close()
		_ = mysql.UserDrop(user.Username, user.Host)
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		defer postgres.Close()
		_ = postgres.UserDrop(user.Username)
	}

	return r.db.Where("id = ?", id).Delete(&biz.DatabaseUser{}).Error
}

func (r databaseUserRepo) DeleteByNames(serverID uint, names []string) error {
	server, err := r.server.Get(serverID)
	if err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		defer mysql.Close()
		users := make([]*biz.DatabaseUser, 0)
		if err = r.db.Where("server_id = ? AND username IN ?", serverID, names).Find(&users).Error; err != nil {
			return err
		}
		for name := range slices.Values(names) {
			host := "localhost"
			for u := range slices.Values(users) {
				if u.Username == name {
					host = u.Host
					break
				}
			}
			_ = mysql.UserDrop(name, host)
		}
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		defer postgres.Close()
		for name := range slices.Values(names) {
			_ = postgres.UserDrop(name)
		}
	}

	return r.db.Where("server_id = ? AND username IN ?", serverID, names).Delete(&biz.DatabaseUser{}).Error
}

func (r databaseUserRepo) fillUser(user *biz.DatabaseUser) {
	server, err := r.server.Get(user.ServerID)
	if err == nil {
		switch server.Type {
		case biz.DatabaseTypeMysql:
			mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
			if err == nil {
				defer mysql.Close()
				privileges, _ := mysql.UserPrivileges(user.Username, user.Host)
				user.Privileges = privileges
			}
			if mysql2, err := db.NewMySQL(user.Username, user.Password, fmt.Sprintf("%s:%d", server.Host, server.Port)); err == nil {
				_ = mysql2.Close()
				user.Status = biz.DatabaseUserStatusValid
			} else {
				user.Status = biz.DatabaseUserStatusInvalid
			}
		case biz.DatabaseTypePostgresql:
			postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
			if err == nil {
				defer postgres.Close()
				privileges, _ := postgres.UserPrivileges(user.Username)
				user.Privileges = privileges
			}
			if postgres2, err := db.NewPostgres(user.Username, user.Password, server.Host, server.Port); err == nil {
				_ = postgres2.Close()
				user.Status = biz.DatabaseUserStatusValid
			} else {
				user.Status = biz.DatabaseUserStatusInvalid
			}
		}
	}
	// 初始化，防止 nil
	if user.Privileges == nil {
		user.Privileges = make([]string, 0)
	}
}
