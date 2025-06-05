package achievement

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllAchievements(c *gin.Context) {
	userID := c.Param("userId")

	var userAchievements []models.UserAchievement
	if err := db.DB.Preload("Achievement").Where("user_id = ?", userID).Find(&userAchievements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
		return
	}

	var responses []models.AchievementResponse
	for _, ua := range userAchievements {
		responses = append(responses, models.AchievementResponse{
			AchievementID: ua.Achievement.ID,
			UserID:        ua.UserID,
			Title:         ua.Achievement.Title,
			Subtitle:      ua.Achievement.Subtitle,
			IsAchieved:    ua.IsAchieved,
		})
	}

	c.JSON(http.StatusOK, responses)
}
