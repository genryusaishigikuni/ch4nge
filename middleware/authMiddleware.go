package middleware

import (
	"fmt"
	"github.com/genryusaishigikuni/ch4nge/jwt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Printf("DEBUG: Authorization header: %s", authHeader)

		if authHeader == "" {
			log.Printf("DEBUG: No authorization header provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("DEBUG: Authorization header doesn't start with 'Bearer '")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("DEBUG: Extracted token: %s", tokenString[:minimum(len(tokenString), 20)]+"...")

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			log.Printf("DEBUG: Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		log.Printf("DEBUG: Token validated successfully for user ID: %d, isAdmin: %t", claims.UserID, claims.IsAdmin)

		c.Set("user_id", claims.UserID)
		c.Set("is_admin", claims.IsAdmin)
		c.Next()
	}
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}
