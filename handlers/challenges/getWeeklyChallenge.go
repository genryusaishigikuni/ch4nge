package challenges

import (
	"log"
	"net/http"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func GetWeeklyChallenge(c *gin.Context) {
	userID := c.Param("userId")

	// Log the request to fetch the weekly challenge for the given user
	log.Printf("Fetching weekly challenge for user ID: %s", userID)

	// Fetch the weekly challenge for the user
	var userChallenge models.UserWeeklyChallenge
	if err := db.DB.Preload("WeeklyChallenge").Where("user_id = ?", userID).First(&userChallenge).Error; err != nil {
		log.Printf("Error fetching weekly challenge for user ID: %s, error: %v", userID, err) // Log error if challenge not found
		c.JSON(http.StatusNotFound, gin.H{"error": "No weekly challenge found"})
		return
	}

	// Log the challenge data that is found
	log.Printf("Weekly challenge found for user ID: %s, Title: %s", userID, userChallenge.WeeklyChallenge.Title)

	// Prepare the response data
	response := models.WeeklyChallengeResponse{
		WeeklyChallengeID: userChallenge.WeeklyChallenge.ID,
		UserID:            userChallenge.UserID,
		Title:             userChallenge.WeeklyChallenge.Title,
		Subtitle:          userChallenge.WeeklyChallenge.Subtitle,
		CurrentValue:      userChallenge.CurrentValue,
		TotalValue:        userChallenge.WeeklyChallenge.TargetValue,
		Points:            userChallenge.WeeklyChallenge.Points,
	}

	// Log the response data before sending it to the client
	log.Printf("Sending weekly challenge response for user ID: %s, Challenge ID: %d", userID, response.WeeklyChallengeID)

	// Send the weekly challenge data in the response
	c.JSON(http.StatusOK, response)
}
