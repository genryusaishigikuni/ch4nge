package admin

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetAllWeeklyChallengesAdmin(c *gin.Context) {
	// Log the request to fetch all weekly challenges
	log.Println("Received request to fetch all weekly challenges")

	// Retrieve all weekly challenges from the database
	var challenges []models.WeeklyChallenge
	if err := database.DB.Find(&challenges).Error; err != nil {
		// Log the error if fetching the weekly challenges fails
		log.Printf("Error fetching weekly challenges: %v", err)

		// Send an error response back to the client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weekly challenges"})
		return
	}

	// Log successful weekly challenge retrieval
	log.Printf("Successfully fetched %d weekly challenges", len(challenges))

	// Send the retrieved weekly challenges back to the client
	c.JSON(http.StatusOK, challenges)
}
