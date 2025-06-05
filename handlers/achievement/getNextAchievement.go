package achievement

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNextAchievement(c *gin.Context) {
	userID := c.Param("userId")

	var userAchievement models.UserAchievement
	if err := db.DB.Preload("Achievement").Where("user_id = ? AND is_achieved = ?", userID, false).First(&userAchievement).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No next achievement found"})
		return
	}

	response := models.AchievementResponse{
		AchievementID: userAchievement.Achievement.ID,
		UserID:        userAchievement.UserID,
		Title:         userAchievement.Achievement.Title,
		Subtitle:      userAchievement.Achievement.Subtitle,
		IsAchieved:    userAchievement.IsAchieved,
	}

	c.JSON(http.StatusOK, response)
}
