package data

import (
	"errors"

	"github.com/go-rat/utils/hash"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/panel"
)

type userRepo struct {
	hasher hash.Hasher
}

func NewUserRepo() biz.UserRepo {
	return &userRepo{
		hasher: hash.NewArgon2id(),
	}
}

func (r *userRepo) CheckPassword(username, password string) (*biz.User, error) {
	user := new(biz.User)
	if err := panel.Orm.Where("username = ?", username).First(user).Error; err != nil {
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
	if err := panel.Orm.First(user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) Save(user *biz.User) error {
	return panel.Orm.Save(user).Error
}
