package repository

import (
	"fmt"
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := NewUserRepository(gormDB)

	test := []struct {
		Name        string
		Username    string
		ExpectedErr error
		MockAct     func()
	}{
		{
			Name:        "Success",
			Username:    "johndoe",
			ExpectedErr: nil,
			MockAct: func() {
				mock.ExpectQuery(config.ExistsTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
		},
		{
			Name:        "Error",
			Username:    "johndoe",
			ExpectedErr: fmt.Errorf("record not found"),
			MockAct: func() {
				mock.ExpectQuery(config.ExistsTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			exists := repo.Exists(tt.Username)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, exists.Error())
			} else {
				assert.NoError(t, exists)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := NewUserRepository(gormDB)

	test := []struct {
		Name        string
		Username    string
		ExpectedErr error
		MockAct     func()
	}{
		{
			Name:        "Success",
			Username:    "johndoe",
			ExpectedErr: nil,
			MockAct: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
		},
		{
			Name:        "Error",
			Username:    "johndoe",
			ExpectedErr: fmt.Errorf("db error"),
			MockAct: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnError(fmt.Errorf("db error"))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			search, searchErr := repo.Search(tt.Username)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, searchErr.Error())
			} else {
				assert.NoError(t, searchErr)
			}
			if search.Username != "" {
				assert.Equal(t, tt.Username, search.Username)
			}
		})
	}
}

func TestSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := NewUserRepository(gormDB)

	test := []struct {
		Name        string
		User        models.User
		ExpectedErr error
		MockAct     func()
	}{
		{
			Name: "Success",
			User: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Phone:    "123456789",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			ExpectedErr: nil,
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.SaveTestQuery).
					WithArgs("1", "John", "Doe", "johndoe", "123456789", "johndoe@example.com", "Password1234").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			Name: "Error",
			User: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Phone:    "123456789",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			ExpectedErr: fmt.Errorf("db error"),
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.SaveTestQuery).
					WithArgs("1", "John", "Doe", "johndoe", "123456789", "johndoe@example.com", "Password1234").
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			save := repo.Save(tt.User)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, save.Error())
			} else {
				assert.NoError(t, save)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := NewUserRepository(gormDB)

	test := []struct {
		Name        string
		Username    string
		Update      models.User
		ExpectedErr error
		MockAct     func()
	}{
		{
			Name:     "Success",
			Username: "johndoe",
			Update: models.User{
				Name:    "Johncito",
				Surname: "Doecito",
				Phone:   "01234567",
				Email:   "johncitodoecito@example.com",
			},
			ExpectedErr: nil,
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.UpdateTestQuery).
					WithArgs("Johncito", "Doecito", "01234567", "johncitodoecito@example.com", "johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			Name:     "Error",
			Username: "johndoe",
			Update: models.User{
				Name:    "Johncito",
				Surname: "Doecito",
				Phone:   "01234567",
				Email:   "johncitodoecito@example.com",
			},
			ExpectedErr: fmt.Errorf("db error"),
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.UpdateTestQuery).
					WithArgs("Johncito", "Doecito", "01234567", "johncitodoecito@example.com", "johndoe").
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
		},
		{
			Name:     "Rows Affected",
			Username: "johndoe",
			Update: models.User{
				Name:    "Johncito",
				Surname: "Doecito",
				Phone:   "01234567",
				Email:   "johncitodoecito@example.com",
			},
			ExpectedErr: fmt.Errorf("no rows affected"),
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.UpdateTestQuery).
					WithArgs("Johncito", "Doecito", "01234567", "johncitodoecito@example.com", "johndoe").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			update := repo.Update(tt.Username, tt.Update)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, update.Error())
			} else {
				assert.NoError(t, update)
			}

		})
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := NewUserRepository(gormDB)

	test := []struct {
		Name        string
		Username    string
		ExpectedErr error
		MockAct     func()
	}{
		{
			Name:        "Success",
			Username:    "johndoe",
			ExpectedErr: nil,
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.DeleteTestQuery).
					WithArgs("johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			Name:        "Error",
			Username:    "johndoe",
			ExpectedErr: fmt.Errorf("db error"),
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.DeleteTestQuery).
					WithArgs("johndoe").
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
		},
		{
			Name:        "Rows Affected",
			Username:    "johndoe",
			ExpectedErr: fmt.Errorf("no rows affected"),
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.DeleteTestQuery).
					WithArgs("johndoe").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			delete := repo.Delete(tt.Username)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, delete.Error())
			} else {
				assert.NoError(t, delete)
			}

		})
	}
}

func TestChangePwd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := NewUserRepository(gormDB)

	test := []struct {
		Name        string
		Username    string
		NewPassword string
		ExpectedErr error
		MockAct     func()
	}{
		{
			Name:        "Success",
			Username:    "johndoe",
			NewPassword: "NewPassword",
			ExpectedErr: nil,
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.ChangePwdTestQuery).
					WithArgs("NewPassword", "johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			Name:        "Error",
			Username:    "johndoe",
			NewPassword: "NewPassword",
			ExpectedErr: fmt.Errorf("db error"),
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.ChangePwdTestQuery).
					WithArgs("NewPassword", "johndoe").
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
		},
		{
			Name:        "Rows Affected",
			Username:    "johndoe",
			NewPassword: "NewPassword",
			ExpectedErr: fmt.Errorf("no rows affected"),
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.ChangePwdTestQuery).
					WithArgs("NewPassword", "johndoe").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			update := repo.ChangePwd(tt.Username, tt.NewPassword)

			if tt.ExpectedErr != nil {
				assert.EqualError(t, tt.ExpectedErr, update.Error())
			} else {
				assert.NoError(t, update)
			}

		})
	}
}
