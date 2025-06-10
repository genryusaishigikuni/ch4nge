package challenges

import (
	"log"
	"net/http"
	"time"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func UpdateUserWeeklyChallengeProgress(c *gin.Context) {
	userID := c.Param("userId")
	challengeID := c.Param("challengeId")

	// Log the request to update progress
	log.Printf("Updating weekly challenge progress for user ID: %s, challenge ID: %s", userID, challengeID)

	var req models.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON for user ID: %s, error: %v", userID, err) // Log the error if request body binding fails
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the user’s challenge data
	var userChallenge models.UserWeeklyChallenge
	if err := db.DB.Where("user_id = ? AND weekly_challenge_id = ?", userID, challengeID).
		First(&userChallenge).Error; err != nil {
		log.Printf("Error fetching user weekly challenge for user ID: %s, challenge ID: %s, error: %v", userID, challengeID, err) // Log error if user challenge is not found
		c.JSON(http.StatusNotFound, gin.H{"error": "User weekly challenge not found"})
		return
	}

	// Log the current progress before update
	log.Printf("Current value before update: %f, adding progress: %f", userChallenge.CurrentValue, req.CurrentValue)

	// Add the progress (can be distance or points based on challenge type)
	userChallenge.CurrentValue += req.CurrentValue

	// Fetch the challenge details
	var challenge models.WeeklyChallenge
	if err := db.DB.First(&challenge, userChallenge.WeeklyChallengeID).Error; err != nil {
		log.Printf("Error fetching challenge details for challenge ID: %s, error: %v", challengeID, err) // Log error if challenge details are not found
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch challenge details"})
		return
	}

	// Check if the challenge is completed
	if userChallenge.CurrentValue >= challenge.TargetValue && !userChallenge.IsCompleted {
		log.Printf("Challenge completed for user ID: %s, challenge ID: %s", userID, challengeID) // Log when the challenge is marked as completed
		userChallenge.IsCompleted = true
		now := time.Now()
		userChallenge.CompletedAt = &now
	}

	// Save the updated challenge progress
	if err := db.DB.Save(&userChallenge).Error; err != nil {
		log.Printf("Error saving updated user challenge progress for user ID: %s, challenge ID: %s, error: %v", userID, challengeID, err) // Log error if save fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update progress"})
		return
	}

	// Log the updated challenge data
	log.Printf("Updated progress for user ID: %s, challenge ID: %s, current value: %f", userID, challengeID, userChallenge.CurrentValue)

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

	// Send the updated progress response
	log.Printf("Sending updated challenge response for user ID: %s, challenge ID: %s", userID, challengeID) // Log sending response
	c.JSON(http.StatusOK, response)
}
