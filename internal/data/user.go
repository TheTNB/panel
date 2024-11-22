package data

import (
	"errors"

	"github.com/go-rat/utils/hash"
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
)

type userRepo struct {
	hasher hash.Hasher
}

func NewUserRepo() biz.UserRepo {
	return do.MustInvoke[biz.UserRepo](injector)
}

func (r *userRepo) Create(username, password string) (*biz.User, error) {
	value, err := r.hasher.Make(password)
	if err != nil {
		return nil, err
	}

	user := &biz.User{
		Username: username,
		Password: value,
	}
	if err = app.Orm.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) CheckPassword(username, password string) (*biz.User, error) {
	user := new(biz.User)
	if err := app.Orm.Where("username = ?", username).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		} else {
			return nil, err
		}
	}

	if !r.hasher.Check(password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	return user, nil
}

func (r *userRepo) Get(id uint) (*biz.User, error) {
	user := new(biz.User)
	if err := app.Orm.First(user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) Save(user *biz.User) error {
	return app.Orm.Save(user).Error
}
