package achievement

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetNextAchievement(c *gin.Context) {
	userID := c.Param("userId")

	// Log the userID parameter to track which user is making the request
	log.Printf("Fetching next achievement for user: %s", userID)

	// Query for the next unachieved achievement for the user
	var userAchievement models.UserAchievement
	if err := db.DB.Preload("Achievement").Where("user_id = ? AND is_achieved = ?", userID, false).First(&userAchievement).Error; err != nil {
		// Log the error if it happens
		log.Printf("Error fetching next achievement for user %s: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "No next achievement found"})
		return
	}

	// Log the achievement found
	log.Printf("Next unachieved achievement found for user %s: %s (ID: %d)", userID, userAchievement.Achievement.Title, userAchievement.Achievement.ID)

	// Prepare the response
	response := models.AchievementResponse{
		AchievementID: userAchievement.Achievement.ID,
		UserID:        userAchievement.UserID,
		Title:         userAchievement.Achievement.Title,
		Subtitle:      userAchievement.Achievement.Subtitle,
		IsAchieved:    userAchievement.IsAchieved,
	}

	// Log the response before sending it
	log.Printf("Returning next achievement: %s (ID: %d) for user %s", userAchievement.Achievement.Title, userAchievement.Achievement.ID, userID)

	// Send the response
	c.JSON(http.StatusOK, response)
}
