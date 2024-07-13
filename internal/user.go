package internal

import "github.com/TheTNB/panel/v2/app/models"

type User interface {
	Create(name, password string) (models.User, error)
	Update(user models.User) (models.User, error)
}
