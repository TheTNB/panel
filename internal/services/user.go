// Package services 用户服务
package services

import (
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/models"
)

type UserImpl struct {
}

func NewUserImpl() *UserImpl {
	return &UserImpl{}
}

func (r *UserImpl) Create(username, password string) (models.User, error) {
	user := models.User{
		Username: username,
		Password: password,
	}
	if err := facades.Orm().Query().Create(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserImpl) Update(user models.User) (models.User, error) {
	if _, err := facades.Orm().Query().Update(&user); err != nil {
		return user, err
	}

	return user, nil
}
