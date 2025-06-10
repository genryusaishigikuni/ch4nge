package admin

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func DeleteAchievement(c *gin.Context) {
	// Extract the achievement ID from the request URL parameter
	achievementID := c.Param("id")

	// Log the attempt to delete the achievement
	log.Printf("Received request to delete achievement with ID: %s", achievementID)

	// Attempt to delete the achievement from the database
	if err := database.DB.Delete(&models.Achievement{}, achievementID).Error; err != nil {
		// Log the error if deletion fails
		log.Printf("Error deleting achievement with ID %s: %v", achievementID, err)

		// Send error response back to the client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete achievement"})
		return
	}

	// Log successful deletion
	log.Printf("Achievement with ID %s deleted successfully", achievementID)

	// Send success response
	c.JSON(http.StatusOK, gin.H{"message": "Achievement deleted successfully"})
}
