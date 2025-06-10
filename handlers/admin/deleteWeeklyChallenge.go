package admin

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func DeleteWeeklyChallenge(c *gin.Context) {
	// Extract the weekly challenge ID from the request URL parameter
	challengeID := c.Param("id")

	// Log the attempt to delete the weekly challenge
	log.Printf("Received request to delete weekly challenge with ID: %s", challengeID)

	// Attempt to delete the weekly challenge from the database
	if err := database.DB.Delete(&models.WeeklyChallenge{}, challengeID).Error; err != nil {
		// Log the error if deletion fails
		log.Printf("Error deleting weekly challenge with ID %s: %v", challengeID, err)

		// Send error response back to the client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete weekly challenge"})
		return
	}

	// Log successful deletion
	log.Printf("Weekly challenge with ID %s deleted successfully", challengeID)

	// Send success response
	c.JSON(http.StatusOK, gin.H{"message": "Weekly challenge deleted successfully"})
}
