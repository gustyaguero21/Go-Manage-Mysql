package router

import (
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/handlers"
	"go-manage-mysql/internal/middleware"
	"go-manage-mysql/internal/repository"
	"go-manage-mysql/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UrlMapping(r *gin.Engine, conn *gorm.DB) {
	api := r.Group(config.BaseURL)

	repo := repository.NewUserRepository(conn)
	service := services.NewUserServices(repo)
	handler := handlers.NewUserHandler(service)

	api.POST("/login", handler.LoginUserHandler)
	api.POST("/create", handler.CreateUserHandler)

	protected := api.Group("/")
	protected.Use(middleware.JWTMiddleware())

	protected.GET("/search", handler.SearchUserHandler)
	protected.PATCH("/update", handler.UpdateUserHandler)
	protected.DELETE("/delete", handler.DeleteUserHandler)
	protected.PATCH("/change-password", handler.ChangePwdHandler)
}
