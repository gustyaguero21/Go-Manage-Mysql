package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(conn *gorm.DB) *gin.Engine {
	router := gin.Default()

	UrlMapping(router, conn)

	return router
}
