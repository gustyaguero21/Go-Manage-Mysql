package repository

import (
	"errors"
	"fmt"
	"go-manage-mysql/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Exists(username string) (bool, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repository) Search(username string) (models.User, error) {
	var user models.User
	result := r.DB.Where("username=?", username).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *Repository) Save(user models.User) error {
	result := r.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) Update(username string, update models.User) error {
	result := r.DB.Where("username=?", username).Select("name", "surname", "phone", "email").Updates(&update)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (r *Repository) Delete(username string) error {
	result := r.DB.Model(&models.User{}).Where("username=?", username).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *Repository) ChangePwd(username string, newPwd string) error {
	result := r.DB.Model(&models.User{}).Where("username = ?", username).Update("password", newPwd)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
