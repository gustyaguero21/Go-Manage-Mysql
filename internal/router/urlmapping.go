package router

import (
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/handlers"
	"go-manage-mysql/internal/repository"
	"go-manage-mysql/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UrlMapping(r *gin.Engine, conn *gorm.DB) {
	api := r.Group(config.BaseURL)

	repo := repository.NewUserRepository(conn)
	service := services.NewUserServices(repo)
	handler := handlers.NewUserHandler(service)

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api.POST("/create", handler.CreateUserHandler)
	api.GET("/search", handler.SearchUserHandler)
	api.PATCH("/update", handler.UpdataUserHandler)
	api.DELETE("/delete", handler.DeleteUserHandler)

}
