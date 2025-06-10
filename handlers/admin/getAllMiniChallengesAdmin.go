package admin

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetAllMiniChallengesAdmin(c *gin.Context) {
	// Log the request to fetch all mini challenges
	log.Println("Received request to fetch all mini challenges")

	// Retrieve all mini challenges from the database
	var challenges []models.MiniChallenge
	if err := database.DB.Find(&challenges).Error; err != nil {
		// Log the error if fetching the mini challenges fails
		log.Printf("Error fetching mini challenges: %v", err)

		// Send an error response back to the client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mini challenges"})
		return
	}

	// Log successful mini challenge retrieval
	log.Printf("Successfully fetched %d mini challenges", len(challenges))

	// Send the retrieved mini challenges back to the client
	c.JSON(http.StatusOK, challenges)
}
