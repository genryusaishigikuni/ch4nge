package achievement

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetAchievementProgress(c *gin.Context) {
	userID := c.Param("userId")

	// Log the userID parameter to track which user is making the request
	log.Printf("Fetching achievement progress for user: %s", userID)

	// Query for the user's achievements
	var userAchievements []models.UserAchievement
	if err := db.DB.Preload("Achievement").Where("user_id = ?", userID).Order("is_achieved DESC, id ASC").Limit(3).Find(&userAchievements).Error; err != nil {
		log.Printf("Error fetching achievements for user %s: %v", userID, err) // Log the error if it happens
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievement progress"})
		return
	}

	// Log the number of achievements found for this user
	log.Printf("Found %d achievements for user %s", len(userAchievements), userID)

	var responses []models.AchievementResponse
	for _, ua := range userAchievements {
		// Log each achievement being processed
		log.Printf("Processing achievement: %s (ID: %d, IsAchieved: %v)", ua.Achievement.Title, ua.Achievement.ID, ua.IsAchieved)
		responses = append(responses, models.AchievementResponse{
			AchievementID: ua.Achievement.ID,
			UserID:        ua.UserID,
			Title:         ua.Achievement.Title,
			Subtitle:      ua.Achievement.Subtitle,
			IsAchieved:    ua.IsAchieved,
		})
	}

	// Log the response data before sending it back to the client
	log.Printf("Returning %d achievements to user %s", len(responses), userID)

	// Send the response
	c.JSON(http.StatusOK, responses)
}
