package router

import (
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UrlMapping(r *gin.Engine) {
	api := r.Group(config.BaseURL)

	_, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
