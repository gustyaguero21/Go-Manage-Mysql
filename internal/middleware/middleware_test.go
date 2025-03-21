package middleware

import (
	"go-manage-mysql/cmd/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func TestJWTMiddleware(t *testing.T) {

	claims := jwt.MapClaims{"username": "testuser"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.GetToken()))

	tests := []struct {
		name         string
		token        string
		expectStatus int
	}{
		{"Valid Token", "Bearer " + tokenString, http.StatusOK},
		{"Missing Token", "", http.StatusUnauthorized},
		{"Invalid Token Format", "InvalidToken", http.StatusUnauthorized},
		{"Invalid Signature", "Bearer invalid.token.signature", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar router de prueba
			r := gin.Default()
			r.Use(JWTMiddleware())
			r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

			// Crear solicitud
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectStatus, w.Code)
			}
		})
	}
}
