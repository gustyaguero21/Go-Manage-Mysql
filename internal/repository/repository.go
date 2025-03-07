package repository

import "go-manage-mysql/internal/models"

type UserRepository interface {
	Search(username string) (models.User, error)
	Save(user models.User) error
	Update(username string, update models.User) error
	Delete(username string) error
	ChangePwd(username string, newPwd string) error
}
