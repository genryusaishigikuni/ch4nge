package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AssignMiniChallengeToUser(c *gin.Context) {
	var req struct {
		UserID          uint `json:"userId" binding:"required"`
		MiniChallengeID uint `json:"miniChallengeId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userChallenge := models.UserMiniChallenge{
		UserID:          req.UserID,
		MiniChallengeID: req.MiniChallengeID,
		AssignedAt:      time.Now(),
	}

	if err := db.DB.Create(&userChallenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign mini challenge"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Mini challenge assigned successfully"})
}
