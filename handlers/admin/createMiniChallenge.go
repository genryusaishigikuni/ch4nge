package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreateMiniChallenge(c *gin.Context) {
	var challenge models.MiniChallenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the mini challenge
	if err := db.DB.Create(&challenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mini challenge"})
		return
	}

	// Assign to all existing users
	if err := assignMiniChallengeToAllUsers(challenge.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Mini challenge created but failed to assign to all users"})
		return
	}

	c.JSON(http.StatusCreated, challenge)
}

func assignMiniChallengeToAllUsers(challengeID uint) error {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		userChallenge := models.UserMiniChallenge{
			UserID:          user.ID,
			MiniChallengeID: challengeID,
			AssignedAt:      time.Now(),
		}
		// Use FirstOrCreate to avoid duplicates
		if err := db.DB.FirstOrCreate(&userChallenge, models.UserMiniChallenge{
			UserID:          user.ID,
			MiniChallengeID: challengeID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
