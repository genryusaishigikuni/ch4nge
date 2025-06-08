package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateWeeklyChallenge(c *gin.Context) {
	challengeID := c.Param("id")

	var challenge models.WeeklyChallenge
	if err := db.DB.First(&challenge, challengeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Weekly challenge not found"})
		return
	}

	var req models.WeeklyChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	challenge.Title = req.Title
	challenge.Subtitle = req.Subtitle
	challenge.Points = req.Points
	challenge.TargetValue = req.TargetValue
	challenge.IsActive = req.IsActive
	challenge.StartDate = req.StartDate
	challenge.EndDate = req.EndDate

	if err := db.DB.Save(&challenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update weekly challenge"})
		return
	}

	c.JSON(http.StatusOK, challenge)
}
