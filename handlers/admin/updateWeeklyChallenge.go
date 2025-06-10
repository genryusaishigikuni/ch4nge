package admin

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func UpdateWeeklyChallenge(c *gin.Context) {
	challengeID := c.Param("id")

	// Log the incoming request for updating the weekly challenge
	log.Printf("Received request to update weekly challenge with ID: %s", challengeID)

	// Retrieve the existing weekly challenge from the database
	var challenge models.WeeklyChallenge
	if err := database.DB.First(&challenge, challengeID).Error; err != nil {
		// Log the error if the challenge is not found
		log.Printf("Weekly challenge with ID %s not found: %v", challengeID, err)

		// Send an error response if the challenge is not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Weekly challenge not found"})
		return
	}

	// Bind the incoming JSON data to the weekly challenge request model
	var req models.WeeklyChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Log the error if binding the JSON fails
		log.Printf("Failed to bind request data for weekly challenge update: %v", err)

		// Send an error response if the binding fails
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the binding of the new data before applying changes
	log.Printf("Updating weekly challenge with new data: Title=%s, Subtitle=%s, Points=%d, TargetValue=%f, IsActive=%v, StartDate=%s, EndDate=%s",
		req.Title, req.Subtitle, req.Points, req.TargetValue, req.IsActive, req.StartDate, req.EndDate)

	// Apply the updates to the challenge
	challenge.Title = req.Title
	challenge.Subtitle = req.Subtitle
	challenge.Points = req.Points
	challenge.TargetValue = req.TargetValue
	challenge.IsActive = req.IsActive
	challenge.StartDate = req.StartDate
	challenge.EndDate = req.EndDate

	// Save the updated weekly challenge to the database
	if err := database.DB.Save(&challenge).Error; err != nil {
		// Log the error if saving the updated challenge fails
		log.Printf("Failed to update weekly challenge with ID %s: %v", challengeID, err)

		// Send an error response if saving the updated challenge fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update weekly challenge"})
		return
	}

	// Log the successful update of the challenge
	log.Printf("Successfully updated weekly challenge with ID: %s", challengeID)

	// Send the updated challenge as the response
	c.JSON(http.StatusOK, challenge)
}
