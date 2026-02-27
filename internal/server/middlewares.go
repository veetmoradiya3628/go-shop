package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/veetmoradiya3628/go-shop/internal/models"
	"github.com/veetmoradiya3628/go-shop/internal/utils"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid Authorization header format")
			c.Abort()
			return
		}
		claims, err := utils.ValidateToken(tokenParts[1], s.config.JWT.Secret)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid token: "+err.Error())
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != string(models.UserRoleAdmin) {
			utils.ForbiddenResponse(c, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
