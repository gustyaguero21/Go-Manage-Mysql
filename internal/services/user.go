package services

import (
	"context"
	"fmt"
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/models"
	"go-manage-mysql/internal/repository"

	"github.com/google/uuid"
	"github.com/gustyaguero21/go-core/pkg/apperror"
	"github.com/gustyaguero21/go-core/pkg/encrypter"
)

type Services struct {
	Repo repository.UserRepository
}

func NewUserServices(repo repository.UserRepository) *Services {
	return &Services{Repo: repo}
}

func (s *Services) CreateUser(ctx context.Context, user models.User) (created models.User, err error) {
	if err := s.exists(user.Username); err != nil {
		return models.User{}, err
	}

	user.ID = uuid.NewString()

	hash, hashErr := encrypter.PasswordEncrypter(user.Password)
	if hashErr != nil {
		return models.User{}, hashErr
	}
	user.Password = string(hash)

	if err := s.Repo.Save(user); err != nil {
		return models.User{}, apperror.AppError(config.ErrCreatingUser, err)
	}

	return user, nil
}

func (s *Services) SearchUser(ctx context.Context, username string) (user models.User, err error) {
	search, searchErr := s.Repo.Search(username)
	if searchErr != nil {
		return models.User{}, apperror.AppError(config.ErrSearchingUser, searchErr)
	}

	return search, nil
}

func (s *Services) UpdateUser(ctx context.Context, username string, update models.User) (err error) {
	if err := s.ensureExists(username); err != nil {
		return err
	}

	if updateErr := s.Repo.Update(username, update); updateErr != nil {
		return apperror.AppError(config.ErrUpdatingUser, updateErr)
	}

	return nil
}

func (s *Services) DeleteUser(ctx context.Context, username string) (err error) {
	if err := s.ensureExists(username); err != nil {
		return err
	}

	if deleteErr := s.Repo.Delete(username); deleteErr != nil {
		return apperror.AppError(config.ErrDeletingUser, deleteErr)
	}
	return nil
}

func (s *Services) ChangeUserPwd(ctx context.Context, username string, newPwd string) (err error) {
	if err := s.ensureExists(username); err != nil {
		return err
	}

	hash, hashErr := encrypter.PasswordEncrypter(newPwd)
	if hashErr != nil {
		return hashErr
	}

	if changeErr := s.Repo.ChangePwd(username, string(hash)); changeErr != nil {
		return apperror.AppError(config.ErrChangingPwd, changeErr)
	}

	return nil
}

func (s *Services) exists(username string) error {
	exists := s.Repo.Exists(username)
	if exists {
		return fmt.Errorf(config.ErrUserAlreadyExists)
	}
	return nil
}

func (s *Services) ensureExists(username string) error {
	exists := s.Repo.Exists(username)
	if !exists {
		return fmt.Errorf(config.ErrUserNotFound)
	}
	return nil
}
