package services

import (
	"context"
	"fmt"
	"go-manage-mysql/internal/models"
	"go-manage-mysql/internal/repository"

	"github.com/google/uuid"
	"github.com/gustyaguero21/go-core/pkg/apperror"
	"github.com/gustyaguero21/go-core/pkg/encrypter"
)

type Services struct {
	Repo repository.Repository
}

func NewUserServices(repo repository.Repository) *Services {
	return &Services{Repo: repo}
}

func (s *Services) CreateUser(ctx context.Context, user models.User) (created models.User, err error) {
	exists, existsErr := s.Repo.Exists(user.Username)
	if existsErr != nil {
		return models.User{}, apperror.AppError("error finding user data", existsErr)
	}
	if exists {
		return models.User{}, fmt.Errorf("user already exists")
	}

	user.ID = uuid.NewString()

	hash, hashErr := encrypter.PasswordEncrypter(user.Password)
	if hashErr != nil {
		return models.User{}, hashErr
	}
	user.Password = string(hash)

	if err := s.Repo.Save(user); err != nil {
		return models.User{}, apperror.AppError("error creating user", err)
	}

	return user, nil
}

func (s *Services) SearchUser(ctx context.Context, username string) (user models.User, err error) {
	search, searchErr := s.Repo.Search(username)
	if searchErr != nil {
		return models.User{}, apperror.AppError("error searching user", searchErr)
	}

	return search, nil
}

func (s *Services) UpdateUser(ctx context.Context, username string, update models.User) (err error) {
	exists, existsErr := s.Repo.Exists(username)
	if existsErr != nil {
		return apperror.AppError("error finding user data", existsErr)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	if updateErr := s.Repo.Update(username, update); updateErr != nil {
		return apperror.AppError("error updating user data", updateErr)
	}

	return nil
}

func (s *Services) DeleteUser(ctx context.Context, username string) (err error) {
	exists, existsErr := s.Repo.Exists(username)
	if existsErr != nil {
		return apperror.AppError("error finding user data", existsErr)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}
	if deleteErr := s.Repo.Delete(username); deleteErr != nil {
		return apperror.AppError("error deleting user data", deleteErr)
	}
	return nil
}

func (s *Services) ChangeUserPwd(ctx context.Context, username string, newPwd string) (err error) {
	exists, existsErr := s.Repo.Exists(username)
	if existsErr != nil {
		return apperror.AppError("error finding user data", existsErr)
	}

	if !exists {
		return fmt.Errorf("user not found")
	}

	if changeErr := s.Repo.ChangePwd(username, newPwd); changeErr != nil {
		return apperror.AppError("error changing user password", changeErr)
	}

	return nil
}
