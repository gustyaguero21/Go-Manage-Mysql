package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWTMiddleware(t *testing.T) {

	jwtSecret = []byte("valid-token")

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(jwtSecret)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Token",
			authHeader:     "Bearer " + tokenString,
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "Missing Token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"status":401,"error":"required token"}`,
		},
		{
			name:           "Invalid Token Format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"status":401,"error":"invalid token format"}`,
		},
		{
			name:           "Invalid Bearer Format",
			authHeader:     "Bearer",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"status":401,"error":"invalid token format"}`,
		},
		{
			name:           "Invalid Token ",
			authHeader:     "Bearer invalid.token.signature",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"status":401,"error":"invalid token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(JWTMiddleware())
			router.GET("/test", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
