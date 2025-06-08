package challenges

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func UpdateUserWeeklyChallengeProgress(c *gin.Context) {
	userID := c.Param("userId")
	challengeID := c.Param("challengeId")

	var req models.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userChallenge models.UserWeeklyChallenge
	if err := db.DB.Where("user_id = ? AND weekly_challenge_id = ?", userID, challengeID).
		First(&userChallenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User weekly challenge not found"})
		return
	}

	userChallenge.CurrentValue = req.CurrentValue

	var challenge models.WeeklyChallenge
	if err := db.DB.First(&challenge, userChallenge.WeeklyChallengeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch challenge details"})
		return
	}

	if userChallenge.CurrentValue >= challenge.TargetValue && !userChallenge.IsCompleted {
		userChallenge.IsCompleted = true
		now := time.Now()
		userChallenge.CompletedAt = &now
	}

	if err := db.DB.Save(&userChallenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update progress"})
		return
	}

	response := models.WeeklyChallengeResponse{
		WeeklyChallengeID: userChallenge.WeeklyChallengeID,
		UserID:            userChallenge.UserID,
		Title:             challenge.Title,
		Subtitle:          challenge.Subtitle,
		CurrentValue:      userChallenge.CurrentValue,
		TotalValue:        challenge.TargetValue,
		Points:            challenge.Points,
	}

	c.JSON(http.StatusOK, response)
}
