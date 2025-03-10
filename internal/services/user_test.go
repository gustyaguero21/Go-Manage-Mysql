package services

import (
	"context"
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/models"
	"go-manage-mysql/internal/repository"
	"go-manage-mysql/internal/utils/testutils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gustyaguero21/go-core/pkg/apperror"
	"github.com/gustyaguero21/go-core/pkg/encrypter"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := NewUserServices(repo)

	test := []struct {
		Name        string
		User        models.User
		ExpectedErr error
		ExistsMock  func()
		MockAct     func()
	}{
		{
			Name:        "Exists Error",
			User:        testutils.OpenMock("../mocks/user.json"),
			ExpectedErr: config.ErrDbError,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnError(config.ErrDbError)
			},
			MockAct: func() {
			},
		},
		{
			Name:        "User already exists",
			User:        testutils.OpenMock("../mocks/user.json"),
			ExpectedErr: apperror.AppError(config.ErrCreatingUser, config.ErrUserAlreadyExists),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "Error creating user",
			User:        testutils.OpenMock("../mocks/user.json"),
			ExpectedErr: apperror.AppError(config.ErrCreatingUser, config.ErrDbError),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.SaveTestQuery).
					WithArgs(sqlmock.AnyArg(), "John", "Doe", "johndoe", "123456789", "johndoe@example.com", sqlmock.AnyArg()).
					WillReturnError(config.ErrDbError)
				mock.ExpectRollback()
			},
		},
		{
			Name:        "Success",
			User:        testutils.OpenMock("../mocks/user.json"),
			ExpectedErr: nil,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.SaveTestQuery).
					WithArgs(sqlmock.AnyArg(), "John", "Doe", "johndoe", "123456789", "johndoe@example.com", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			create, createErr := service.CreateUser(ctx, tt.User)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, createErr.Error())
			} else {
				assert.NoError(t, createErr)
			}

			if create.ID != "" {
				assert.Equal(t, tt.User.Username, create.Username)
			}
		})
	}
}

func TestSearchUser(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := NewUserServices(repo)

	test := []struct {
		Name         string
		Username     string
		ExpectedUser models.User
		ExpectedErr  error
		MockAct      func()
	}{
		{
			Name:         "Error searching user",
			Username:     "johndoe",
			ExpectedUser: testutils.OpenMock("../mocks/user.json"),
			ExpectedErr:  apperror.AppError(config.ErrSearchingUser, config.ErrUserNotFound),
			MockAct: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
		},
		{
			Name:         "Success",
			Username:     "johndoe",
			ExpectedUser: testutils.OpenMock("../mocks/user.json"),
			ExpectedErr:  nil,
			MockAct: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "phone", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "123456789", "johndoe@example.com", "hashedpassword"))

			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			search, searchErr := service.SearchUser(ctx, tt.Username)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, searchErr.Error())
			} else {
				assert.NoError(t, searchErr)
			}

			if search.ID != "" {
				assert.Equal(t, tt.ExpectedUser.Username, search.Username)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := NewUserServices(repo)

	test := []struct {
		Name        string
		Username    string
		Update      models.User
		ExpectedErr error
		ExistsMock  func()
		MockAct     func()
	}{
		{
			Name:        "Exists Error",
			Username:    "johndoe",
			Update:      testutils.OpenMock("../mocks/update_user.json"),
			ExpectedErr: apperror.AppError(config.ErrUpdatingUser, config.ErrDbError),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnError(apperror.AppError(config.ErrUpdatingUser, config.ErrDbError))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "User not found",
			Username:    "johndoe",
			Update:      testutils.OpenMock("../mocks/update_user.json"),
			ExpectedErr: apperror.AppError(config.ErrUpdatingUser, config.ErrUserNotFound),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "Error updating user",
			Username:    "johndoe",
			Update:      testutils.OpenMock("../mocks/update_user.json"),
			ExpectedErr: apperror.AppError(config.ErrUpdatingUser, config.ErrNoNewData),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.UpdateTestQuery).
					WithArgs("Johncito", "Doecito", "23456789", "johncitodoecito@example.com", "johndoe").
					WillReturnError(apperror.AppError(config.ErrUpdatingUser, config.ErrNoNewData))
				mock.ExpectRollback()
			},
		},
		{
			Name:        "Success",
			Username:    "johndoe",
			Update:      testutils.OpenMock("../mocks/update_user.json"),
			ExpectedErr: nil,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.UpdateTestQuery).
					WithArgs("Johncito", "Doecito", "23456789", "johncitodoecito@example.com", "johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			update := service.UpdateUser(ctx, tt.Username, tt.Update)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, update.Error())
			} else {
				assert.NoError(t, update)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := NewUserServices(repo)

	test := []struct {
		Name        string
		Username    string
		ExpectedErr error
		ExistsMock  func()
		MockAct     func()
	}{
		{
			Name:        "Exists Error",
			Username:    "johndoe",
			ExpectedErr: apperror.AppError(config.ErrDeletingUser, config.ErrDbError),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnError(apperror.AppError(config.ErrDeletingUser, config.ErrDbError))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "User not found",
			Username:    "johndoe",
			ExpectedErr: apperror.AppError(config.ErrDeletingUser, config.ErrUserNotFound),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "Error deleting user",
			Username:    "johndoe",
			ExpectedErr: apperror.AppError(config.ErrDeletingUser, config.ErrDbError),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.DeleteTestQuery).
					WithArgs("johndoe").
					WillReturnError(config.ErrDbError)
				mock.ExpectRollback()
			},
		},
		{
			Name:        "Success",
			Username:    "johndoe",
			ExpectedErr: nil,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.DeleteTestQuery).
					WithArgs("johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			delete := service.DeleteUser(ctx, tt.Username)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, delete.Error())
			} else {
				assert.NoError(t, delete)
			}
		})
	}
}

func TestChangePwd(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := NewUserServices(repo)

	test := []struct {
		Name        string
		Username    string
		NewPwd      string
		ExpectedErr error
		ExistsMock  func()
		MockAct     func()
	}{
		{
			Name:        "Exists Error",
			Username:    "johndoe",
			ExpectedErr: apperror.AppError(config.ErrChangingPwd, config.ErrDbError),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnError(apperror.AppError(config.ErrChangingPwd, config.ErrDbError))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "User not found",
			Username:    "johndoe",
			NewPwd:      "NewPassword1234",
			ExpectedErr: apperror.AppError(config.ErrChangingPwd, config.ErrUserNotFound),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "Error changing user password",
			Username:    "johndoe",
			NewPwd:      "NewPassword1234",
			ExpectedErr: apperror.AppError(config.ErrChangingPwd, config.ErrDbError),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.ChangePwdTestQuery).
					WithArgs(sqlmock.AnyArg(), "johndoe").
					WillReturnError(config.ErrDbError)
				mock.ExpectRollback()
			},
		},
		{
			Name:        "Success",
			Username:    "johndoe",
			NewPwd:      "NewPassword1234",
			ExpectedErr: nil,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.ChangePwdTestQuery).
					WithArgs(sqlmock.AnyArg(), "johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			change := service.ChangeUserPwd(ctx, tt.Username, tt.NewPwd)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, change.Error())
			} else {
				assert.NoError(t, change)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := NewUserServices(repo)

	tests := []struct {
		Name        string
		Username    string
		Password    string
		ExpectedErr error
		ExistsMock  func()
		MockAct     func()
	}{
		{
			Name:        "Exists Error",
			Username:    "johndoe",
			Password:    "Password1234",
			ExpectedErr: apperror.AppError(config.ErrLoginUser, config.ErrDbError),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnError(apperror.AppError(config.ErrLoginUser, config.ErrDbError))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "User not found",
			Username:    "johndoe",
			Password:    "Password1234",
			ExpectedErr: apperror.AppError(config.ErrLoginUser, config.ErrUserNotFound),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "Passwords doesn't match",
			Username:    "johndoe",
			Password:    "Password12",
			ExpectedErr: apperror.AppError(config.ErrLoginUser, config.ErrPwdMatching),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
				hashedPwd, _ := encrypter.PasswordEncrypter("Password1234")

				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).
						AddRow(1, hashedPwd))
			},
		},
		{
			Name:        "Error login user",
			Username:    "johndoe",
			Password:    "Password1234",
			ExpectedErr: apperror.AppError(config.ErrSearchingUser, config.ErrDbError),
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnError(config.ErrDbError)
			},
		},

		{
			Name:        "Success",
			Username:    "johndoe",
			Password:    "Password1234",
			ExpectedErr: nil,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			MockAct: func() {
				hashedPwd, _ := encrypter.PasswordEncrypter("Password1234")

				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).
						AddRow(1, hashedPwd))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			err := service.LoginUser(ctx, tt.Username, tt.Password)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
