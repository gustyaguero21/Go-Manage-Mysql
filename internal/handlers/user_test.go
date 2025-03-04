package handlers

import (
	"bytes"
	"fmt"
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/repository"
	"go-manage-mysql/internal/services"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
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
			Name: "Success",
			Body: `{
				"id":"1",
				"name": "John",
				"surname": "Doe",
				"username": "johndoe",
				"phone":"23456789",
				"email": "johndoe@example.com",
				"password": "Password1234"
			}`,
			ExpectedCode: http.StatusCreated,
			ExistsMock: func() {
				mock.ExpectQuery(config.ExistsTestQuery).
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
			Body:         `{"id": "1", "name": "John", "surname": "Doe", "username": "johndoe","phone":"123456789", "email": "johndoe@example.com", "password": }`,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Validate Error",
			Body:         `{"id": "1", "name": "John", "surname": "Doe", "username": "johndoe","phone":"123456789"}`,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name: "Error",
			Body: `{
				"id": "1",
				"name": "John",
				"surname": "Doe",
				"username": "johndoe",
				"phone":"23456789",
				"email": "johndoe@example.com",
				"password": "Password1234"
			}`,
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
			Name:     "Invalid Query Param",
			Username: "",
			Body: `{
				"surname": "Doecito",
				"phone":"23456789",
				"email": "johncitodoecito@example.com"
			}`,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Invalid JSON Body",
			Username:     "johndoe",
			Body:         `{"name": "Johncito", "surname": "Doecito", "phone": "23456789", "email": "johncitodoecito@example.com"`,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:     "Validation Error",
			Username: "johndoe",
			Body: `{
				"surname": "Doecito",
				"phone":"23456789",
				"email": "johncitodoecito@"
			}`,
			ExpectedCode: http.StatusBadRequest,
			ExistsMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:     "Success",
			Username: "johndoe",
			Body: `{
				"name": "Johncito",
				"surname": "Doecito",
				"phone":"23456789",
				"email": "johncitodoecito@example.com"
			}`,
			ExpectedCode: http.StatusOK,
			ExistsMock: func() {
				mock.ExpectQuery(config.ExistsTestQuery).
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
			Name:     "Error",
			Username: "johndoe",
			Body: `{
				"id":"1",
				"name": "Johncito",
				"surname": "Doecito",
				"phone":"23456789",
				"email": "johncitodoecito@example.com"
			}`,
			ExpectedCode: http.StatusInternalServerError,
			ExistsMock: func() {
				mock.ExpectQuery(config.ExistsTestQuery).
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
			Name:         "Invalid Query Param",
			Username:     "",
			ExpectedCode: http.StatusBadRequest,
			ExistsMock: func() {
			},
			MockAct: func() {
			},
		},
		{
			Name:         "Success",
			Username:     "johndoe",
			ExpectedCode: http.StatusOK,
			ExistsMock: func() {
				mock.ExpectQuery(config.ExistsTestQuery).
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
			Name:         "Error",
			Username:     "johndoe",
			ExpectedCode: http.StatusInternalServerError,
			ExistsMock: func() {
				mock.ExpectQuery(config.ExistsTestQuery).
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
		Username     string
		Password     string
		ExpectedCode int
		ExistsMock   func()
		MockAct      func()
	}{
		{
			Name:         "Invalid Query Param",
			Username:     "",
			Password:     "NewPassword1234",
			ExpectedCode: http.StatusBadRequest,
			ExistsMock: func() {
			},
			MockAct: func() {
			},
		},
		{
			Name:         "Success",
			Username:     "johndoe",
			Password:     "NewPassword1234",
			ExpectedCode: http.StatusOK,
			ExistsMock: func() {
				mock.ExpectQuery(config.ExistsTestQuery).
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
			Name:         "Error",
			Username:     "johndoe",
			Password:     "NewPassword1234",
			ExpectedCode: http.StatusInternalServerError,
			ExistsMock: func() {
				mock.ExpectQuery(config.ExistsTestQuery).
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

			url := fmt.Sprintf("/change-password?username=%s&new_password=%s", tt.Username, tt.Password)

			req, _ := http.NewRequest(http.MethodPatch, url, nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}
