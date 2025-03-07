package services

import (
	"context"
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/models"
	"go-manage-mysql/internal/repository"

	"github.com/google/uuid"
	"github.com/gustyaguero21/go-core/pkg/apperror"
	"github.com/gustyaguero21/go-core/pkg/encrypter"
	"gorm.io/gorm"
)

type Services struct {
	Repo repository.UserRepository
}

func NewUserServices(repo repository.UserRepository) *Services {
	return &Services{Repo: repo}
}

func (s *Services) CreateUser(ctx context.Context, user models.User) (created models.User, err error) {
	exist, existErr := s.exists(user.Username)
	if existErr != nil {
		return models.User{}, existErr
	}

	if exist {
		return models.User{}, apperror.AppError(config.ErrUserAlreadyExists, nil)
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
	exist, existErr := s.exists(username)
	if existErr != nil {
		return existErr
	}

	if !exist {
		return apperror.AppError(config.ErrUserNotFound, nil)
	}

	if updateErr := s.Repo.Update(username, update); updateErr != nil {
		return apperror.AppError(config.ErrUpdatingUser, updateErr)
	}

	return nil
}

func (s *Services) DeleteUser(ctx context.Context, username string) (err error) {
	exist, existErr := s.exists(username)
	if existErr != nil {
		return existErr
	}

	if !exist {
		return apperror.AppError(config.ErrUserNotFound, nil)
	}

	if deleteErr := s.Repo.Delete(username); deleteErr != nil {
		return apperror.AppError(config.ErrDeletingUser, deleteErr)
	}
	return nil
}

func (s *Services) ChangeUserPwd(ctx context.Context, username string, newPwd string) (err error) {
	exist, existErr := s.exists(username)
	if existErr != nil {
		return existErr
	}

	if !exist {
		return apperror.AppError(config.ErrUserNotFound, nil)
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

func (s *Services) LoginUser(ctx context.Context, username, password string) error {
	exist, existErr := s.exists(username)
	if existErr != nil {
		return existErr
	}

	if !exist {
		return apperror.AppError(config.ErrUserNotFound, gorm.ErrRecordNotFound)
	}

	search, searchErr := s.Repo.Search(username)
	if searchErr != nil {
		return apperror.AppError(config.ErrSearchingUser, searchErr)
	}

	if !encrypter.PasswordDecrypter([]byte(search.Password), password) {
		return apperror.AppError(config.ErrPwdMatching, nil)
	}
	return nil
}

func (s *Services) exists(username string) (bool, error) {
	search, searchErr := s.Repo.Search(username)
	if searchErr != nil {
		if searchErr == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, searchErr
	}
	return search.ID != "", nil
}
