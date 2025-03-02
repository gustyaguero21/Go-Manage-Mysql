package router

import (
	"go-manage-mysql/cmd/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UrlMapping(r *gin.Engine, conn *gorm.DB) {
	api := r.Group(config.BaseURL)

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
