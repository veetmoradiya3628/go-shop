package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"time"

	"github.com/veetmoradiya3628/go-shop/internal/config"
	"github.com/veetmoradiya3628/go-shop/internal/models"
	"github.com/veetmoradiya3628/go-shop/internal/utils"
)

func performRequest(router *gin.Engine, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestAuthMiddleware(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: "testsecret", ExpiresIn: time.Hour, RefreshTokenExpires: time.Hour},
	}
	s := &Server{
		config: cfg,
	}
	router := gin.New()
	router.Use(s.authMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	// Test missing Authorization header
	w := performRequest(router, "GET", "/test", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test invalid Authorization header format
	w = performRequest(router, "GET", "/test", map[string]string{
		"Authorization": "InvalidFormat",
	})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test valid token by generating one from the utils package
	token, _, err := utils.GenerateTokenPair(&cfg.JWT, 1, "user@example.com", string(models.UserRoleCustomer))
	assert.NoError(t, err)
	w = performRequest(router, "GET", "/test", map[string]string{
		"Authorization": "Bearer " + token,
	})
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestAdminMiddleware(t *testing.T) {
	// helper to create a router with a preset role
	makeRouter := func(role string) *gin.Engine {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_role", role)
			c.Next()
		})
		s := &Server{}
		r.Use(s.adminMiddleware())
		r.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})
		return r
	}

	// non-admin router
	routerNo := makeRouter(string(models.UserRoleCustomer))
	w := performRequest(routerNo, "GET", "/admin", nil)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// admin router
	routerAdmin := makeRouter(string(models.UserRoleAdmin))
	w = performRequest(routerAdmin, "GET", "/admin", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}
