package handlers

import (
	"bytes"
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/mocks"
	"go-manage-mysql/internal/repository"
	"go-manage-mysql/internal/services"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/gustyaguero21/go-core/pkg/encrypter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCreateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := services.NewUserServices(repo)
	handler := NewUserHandler(service)

	r := gin.Default()
	r.POST("/create", handler.CreateUserHandler)

	test := []struct {
		Name         string
		Body         string
		ExpectedCode int
		ExistsMock   func()
		MockAct      func()
	}{
		{
			Name:         "Success",
			Body:         mocks.CreateUser,
			ExpectedCode: http.StatusCreated,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.SaveTestQuery).
					WithArgs().
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			Name:         "Invalid JSON",
			Body:         mocks.InvalidJSON,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Validate Error",
			Body:         mocks.ValidateError,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Error",
			Body:         mocks.CreateUser,
			ExpectedCode: http.StatusInternalServerError,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			req, _ := http.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(tt.Body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}

func TestSearchUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := services.NewUserServices(repo)
	handler := NewUserHandler(service)

	r := gin.Default()
	r.GET("/search", handler.SearchUserHandler)

	tests := []struct {
		Name         string
		Username     string
		ExpectedCode int
		MockAct      func()
	}{
		{
			Name:         "Success",
			Username:     "johndoe",
			ExpectedCode: http.StatusOK,
			MockAct: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
		},
		{
			Name:         "Invalid Query Params",
			Username:     "",
			ExpectedCode: http.StatusBadRequest,
			MockAct: func() {
			},
		},
		{
			Name:         "Error",
			Username:     "johndoe",
			ExpectedCode: http.StatusInternalServerError,
			MockAct: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			url := "/search?username=" + tt.Username
			req, _ := http.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}

func TestUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := services.NewUserServices(repo)
	handler := NewUserHandler(service)

	r := gin.Default()
	r.PATCH("/update", handler.UpdateUserHandler)

	tests := []struct {
		Name         string
		Username     string
		Body         string
		ExpectedCode int
		ExistsMock   func()
		MockAct      func()
	}{
		{
			Name:         "Success",
			Username:     "johndoe",
			Body:         mocks.UpdateUser,
			ExpectedCode: http.StatusOK,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "phone", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "1234567890", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.UpdateTestQuery).
					WithArgs("Johncito", "Doecito", "23456789", "johncitodoecito@example.com", "johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			Name:         "Invalid Query Param",
			Username:     "",
			Body:         mocks.InvalidQueryParam,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Invalid JSON",
			Username:     "johndoe",
			Body:         mocks.InvalidJSON,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Validate Error",
			Username:     "johndoe",
			Body:         mocks.ValidateError,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Error",
			Username:     "johndoe",
			Body:         mocks.UpdateUser,
			ExpectedCode: http.StatusInternalServerError,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "phone", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "1234567890", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			url := "/update?username=" + tt.Username
			req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBufferString(tt.Body))

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}

func TestDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := services.NewUserServices(repo)
	handler := NewUserHandler(service)

	r := gin.Default()
	r.DELETE("/delete", handler.DeleteUserHandler)

	tests := []struct {
		Name         string
		Username     string
		ExpectedCode int
		ExistsMock   func()
		MockAct      func()
	}{
		{
			Name:         "Success",
			Username:     "johndoe",
			ExpectedCode: http.StatusOK,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.DeleteTestQuery).
					WithArgs("johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			Name:         "Invalid Query Param",
			Username:     "",
			ExpectedCode: http.StatusBadRequest,
			ExistsMock: func() {
			},
			MockAct: func() {
			},
		},
		{
			Name:         "Error",
			Username:     "johndoe",
			ExpectedCode: http.StatusInternalServerError,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			url := "/delete?username=" + tt.Username

			req, _ := http.NewRequest(http.MethodDelete, url, nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}

func TestChangePwdHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := services.NewUserServices(repo)
	handler := NewUserHandler(service)

	r := gin.Default()
	r.PATCH("/change-password", handler.ChangePwdHandler)

	test := []struct {
		Name         string
		Body         string
		ExpectedCode int
		ExistsMock   func()
		MockAct      func()
	}{

		{
			Name:         "Success",
			Body:         mocks.ChangePwd,
			ExpectedCode: http.StatusOK,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectBegin()
				mock.ExpectExec(config.ChangePwdTestQuery).
					WithArgs(sqlmock.AnyArg(), "johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			Name:         "Invalid body",
			Body:         mocks.InvalidBody,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock: func() {
			},
			MockAct: func() {
			},
		},
		{
			Name:         "Invalid password format",
			Body:         mocks.InvalidPasswordFormat,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock: func() {
			},
			MockAct: func() {
			},
		},
		{
			Name:         "Error",
			Body:         mocks.ChangePwd,
			ExpectedCode: http.StatusInternalServerError,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			req, _ := http.NewRequest(http.MethodPatch, "/change-password", bytes.NewBufferString(tt.Body))

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}

func TestLoginUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gormDB, gormErr := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if gormErr != nil {
		t.Fatal(gormErr)
	}

	repo := repository.NewUserRepository(gormDB)
	service := services.NewUserServices(repo)
	handler := NewUserHandler(service)

	r := gin.Default()
	r.POST("/login", handler.LoginUserHandler)

	test := []struct {
		Name         string
		Body         string
		ExpectedCode int
		ExistsMock   func()
		MockAct      func()
	}{
		{
			Name: "Success",
			Body: `{
				"username": "johndoe",
				"password":"Password1234"
			}`,
			ExpectedCode: http.StatusOK,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
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
			Name: "Invalid JSON",
			Body: `{
				"username": "johndoe",
				"password":"Password1234",
			}`,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock: func() {
			},
			MockAct: func() {
			},
		},
		{
			Name: "Invalid Query Param",
			Body: `{
				"username": "johndoe",
				"password":""
			}`,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock: func() {
			},
			MockAct: func() {
			},
		},
		{
			Name: "User Not Found",
			Body: `{
				"username": "nonexistent",
				"password":"Password1234"
			}`,
			ExpectedCode: http.StatusNotFound,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("nonexistent", 1).
					WillReturnError(config.ErrRecordNotFound)
			},
			MockAct: func() {},
		},

		{
			Name: "Unauthorized",
			Body: `{
				"username": "johndoe",
				"password":"Password1234"
			}`,
			ExpectedCode: http.StatusUnauthorized,
			ExistsMock: func() {
				mock.ExpectQuery(config.SearchTestQuery).
					WithArgs("johndoe", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.ExistsMock()
			tt.MockAct()

			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.Body))

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}
