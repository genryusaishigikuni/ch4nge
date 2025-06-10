package admin

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func DeleteMiniChallenge(c *gin.Context) {
	// Extract the mini challenge ID from the request URL parameter
	challengeID := c.Param("id")

	// Log the attempt to delete the mini challenge
	log.Printf("Received request to delete mini challenge with ID: %s", challengeID)

	// Attempt to delete the mini challenge from the database
	if err := database.DB.Delete(&models.MiniChallenge{}, challengeID).Error; err != nil {
		// Log the error if deletion fails
		log.Printf("Error deleting mini challenge with ID %s: %v", challengeID, err)

		// Send error response back to the client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete mini challenge"})
		return
	}

	// Log successful deletion
	log.Printf("Mini challenge with ID %s deleted successfully", challengeID)

	// Send success response
	c.JSON(http.StatusOK, gin.H{"message": "Mini challenge deleted successfully"})
}
