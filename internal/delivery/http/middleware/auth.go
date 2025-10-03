package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/islamyakin/otel-propagation-monorepo/internal/entity"
	"github.com/islamyakin/otel-propagation-monorepo/internal/usecase"
)

const (
	AuthUserID   = "user_id"
	AuthUsername = "username"
	AuthRole     = "role"
)

func JWTAuth(authUseCase usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		// Verify token
		claims, err := authUseCase.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(AuthUserID, claims.UserID)
		c.Set(AuthUsername, claims.Username)
		c.Set(AuthRole, claims.Role)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(AuthRole)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No role found in context"})
			c.Abort()
			return
		}

		userRole, ok := role.(entity.Role)
		if !ok || userRole != entity.AdminRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Helper functions to get user info from context
func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get(AuthUserID)
	if !exists {
		return 0, false
	}

	id, ok := userID.(int)
	return id, ok
}

func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(AuthUsername)
	if !exists {
		return "", false
	}

	name, ok := username.(string)
	return name, ok
}

func GetUserRole(c *gin.Context) (entity.Role, bool) {
	role, exists := c.Get(AuthRole)
	if !exists {
		return "", false
	}

	userRole, ok := role.(entity.Role)
	return userRole, ok
}

func IsAdmin(c *gin.Context) bool {
	role, ok := GetUserRole(c)
	return ok && role == entity.AdminRole
}
