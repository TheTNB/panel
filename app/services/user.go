package services

import (
	"github.com/goravel/framework/facades"

	"panel/app/models"
)

type User interface {
	Create(name, password string) (models.User, error)
	Update(user models.User) (models.User, error)
}

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
