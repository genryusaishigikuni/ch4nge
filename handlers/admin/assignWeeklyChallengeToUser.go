package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AssignWeeklyChallengeToUser(c *gin.Context) {
	var req struct {
		UserID            uint `json:"userId" binding:"required"`
		WeeklyChallengeID uint `json:"weeklyChallengeId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userChallenge := models.UserWeeklyChallenge{
		UserID:            req.UserID,
		WeeklyChallengeID: req.WeeklyChallengeID,
		AssignedAt:        time.Now(),
	}

	if err := db.DB.Create(&userChallenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign weekly challenge"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Weekly challenge assigned successfully"})
}
