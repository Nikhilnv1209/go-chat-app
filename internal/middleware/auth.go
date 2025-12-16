package middleware

import (
	"strings"

	"chat-app/internal/errors"
	"chat-app/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT tokens and sets userID in context.
// This follows the specification in specs/03_Technical_Specification.md Section 6.
func AuthMiddleware(jwtService jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": errors.ErrUnauthorized})
			return
		}

		// Extract "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": errors.ErrUnauthorized})
			return
		}

		tokenString := parts[1]
		userID, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": errors.ErrUnauthorized})
			return
		}

		// Set userID in context for downstream handlers
		c.Set("userID", userID)
		c.Next()
	}
}

// GetUserIDFromContext extracts the userID from gin context.
// Returns uuid.Nil if not found or invalid.
func GetUserIDFromContext(c *gin.Context) uuid.UUID {
	val, exists := c.Get("userID")
	if !exists {
		return uuid.Nil
	}
	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}
