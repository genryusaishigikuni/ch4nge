package admin

import (
	"log"
	"net/http"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func AssignAchievementToUser(c *gin.Context) {
	// Log the incoming request
	log.Println("Received request to assign achievement to user")

	// Define the request structure
	var req struct {
		UserID        uint `json:"userId" binding:"required"`
		AchievementID uint `json:"achievementId" binding:"required"`
	}

	// Bind the incoming JSON request to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Request successfully bound. UserID: %d, AchievementID: %d", req.UserID, req.AchievementID)

	// Create the UserAchievement record
	userAchievement := models.UserAchievement{
		UserID:        req.UserID,
		AchievementID: req.AchievementID,
	}

	// Save to the database
	log.Printf("Attempting to assign achievement %d to user %d", req.AchievementID, req.UserID)
	if err := db.DB.Create(&userAchievement).Error; err != nil {
		log.Printf("Error assigning achievement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign achievement"})
		return
	}

	// Success response
	log.Printf("Achievement %d successfully assigned to user %d", req.AchievementID, req.UserID)
	c.JSON(http.StatusCreated, gin.H{"message": "Achievement assigned successfully"})
}
