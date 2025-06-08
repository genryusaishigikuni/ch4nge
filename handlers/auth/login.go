package auth

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"log"

	"github.com/genryusaishigikuni/ch4nge/jwt"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		log.Printf("DEBUG LOGIN: Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Printf("DEBUG LOGIN: Generated token: %s", token)
	log.Printf("DEBUG LOGIN: Token length: %d", len(token))
	log.Printf("DEBUG LOGIN: First 50 chars: %s", token[:minimum(50, len(token))])

	claims, validateErr := jwt.ValidateToken(token)
	if validateErr != nil {
		log.Printf("DEBUG LOGIN: Token validation failed immediately after generation: %v", validateErr)
	} else {
		log.Printf("DEBUG LOGIN: Token validated successfully after generation - UserID: %d", claims.UserID)
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"is_admin": user.IsAdmin,
		},
	})
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}
