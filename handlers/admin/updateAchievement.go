package admin

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func UpdateAchievement(c *gin.Context) {
	achievementID := c.Param("id")

	// Log the request to update the achievement
	log.Printf("Received request to update achievement with ID: %s", achievementID)

	// Retrieve the existing achievement from the database
	var achievement models.Achievement
	if err := database.DB.First(&achievement, achievementID).Error; err != nil {
		// Log the error if the achievement is not found
		log.Printf("Achievement with ID %s not found: %v", achievementID, err)

		// Send an error response back to the client
		c.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}

	// Bind the incoming JSON to the achievement model
	if err := c.ShouldBindJSON(&achievement); err != nil {
		// Log the error if binding the JSON fails
		log.Printf("Failed to bind request data for achievement update: %v", err)

		// Send an error response if JSON binding fails
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated achievement to the database
	if err := database.DB.Save(&achievement).Error; err != nil {
		// Log the error if saving the updated achievement fails
		log.Printf("Failed to update achievement with ID %s: %v", achievementID, err)

		// Send an error response back to the client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievement"})
		return
	}

	// Log the successful achievement update
	log.Printf("Successfully updated achievement with ID: %s", achievementID)

	// Send the updated achievement as the response
	c.JSON(http.StatusOK, achievement)
}
