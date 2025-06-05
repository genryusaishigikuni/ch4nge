package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AssignAchievementToUser(c *gin.Context) {
	var req struct {
		UserID        uint `json:"userId" binding:"required"`
		AchievementID uint `json:"achievementId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAchievement := models.UserAchievement{
		UserID:        req.UserID,
		AchievementID: req.AchievementID,
	}

	if err := db.DB.Create(&userAchievement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign achievement"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Achievement assigned successfully"})
}
