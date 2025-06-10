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

	// Fetch the user’s challenge data
	var userChallenge models.UserWeeklyChallenge
	if err := db.DB.Where("user_id = ? AND weekly_challenge_id = ?", userID, challengeID).
		First(&userChallenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User weekly challenge not found"})
		return
	}

	// Add the progress (can be distance or points based on challenge type)
	userChallenge.CurrentValue += req.CurrentValue

	// Check if the challenge is completed
	var challenge models.WeeklyChallenge
	if err := db.DB.First(&challenge, userChallenge.WeeklyChallengeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch challenge details"})
		return
	}

	// Mark the challenge as completed if the current value exceeds the target
	if userChallenge.CurrentValue >= challenge.TargetValue && !userChallenge.IsCompleted {
		userChallenge.IsCompleted = true
		now := time.Now()
		userChallenge.CompletedAt = &now
	}

	// Save the updated challenge progress
	if err := db.DB.Save(&userChallenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update progress"})
		return
	}

	// Return the updated user challenge data
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
