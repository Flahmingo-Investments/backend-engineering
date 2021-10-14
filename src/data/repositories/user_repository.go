package repositories

import (
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/data/models"

	uuid "github.com/satori/go.uuid"
)

type UserManager interface {
	Create(*models.User) bool
	Update(*models.User) bool
	Get(*models.User) *models.User
}

type UserRepository struct {
}

func (u *UserRepository) Create(user *models.User) bool {
	user.ID = uuid.NewV4().String()
	result := commons.DB.Create(user)
	return result.Error == nil
}

func (u *UserRepository) Update(user *models.User) bool {
	result := commons.DB.Save(user)
	return result.Error == nil && result.RowsAffected > 0
}

func (u *UserRepository) Get(search *models.User) *models.User {
	var user *models.User
	result := commons.DB.Where(search).Find(&user)
	if result.RowsAffected > 0 {
		return user
	}
	return nil
}
