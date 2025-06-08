package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateAchievement(c *gin.Context) {
	var achievement models.Achievement
	if err := c.ShouldBindJSON(&achievement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the achievement
	if err := db.DB.Create(&achievement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create achievement"})
		return
	}

	if err := assignAchievementToAllUsers(achievement.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Achievement created but failed to assign to all users"})
		return
	}

	c.JSON(http.StatusCreated, achievement)
}

func assignAchievementToAllUsers(achievementID uint) error {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		userAchievement := models.UserAchievement{
			UserID:        user.ID,
			AchievementID: achievementID,
		}
		if err := db.DB.FirstOrCreate(&userAchievement, models.UserAchievement{
			UserID:        user.ID,
			AchievementID: achievementID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
