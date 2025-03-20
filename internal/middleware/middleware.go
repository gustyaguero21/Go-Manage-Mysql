package middleware

import (
	"go-manage-mysql/cmd/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gustyaguero21/go-core/pkg/web"
)

var jwtSecret = []byte(config.GetToken())

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			web.NewError(ctx, http.StatusUnauthorized, "required token")
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			web.NewError(ctx, http.StatusUnauthorized, "invalid token format")
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			web.NewError(ctx, http.StatusUnauthorized, "invalid token")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
