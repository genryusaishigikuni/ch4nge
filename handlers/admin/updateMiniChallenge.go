package admin

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func UpdateMiniChallenge(c *gin.Context) {
	challengeID := c.Param("id")

	// Log the request to update the mini challenge
	log.Printf("Received request to update mini challenge with ID: %s", challengeID)

	// Retrieve the existing mini challenge from the database
	var challenge models.MiniChallenge
	if err := database.DB.First(&challenge, challengeID).Error; err != nil {
		// Log the error if the mini challenge is not found
		log.Printf("Mini challenge with ID %s not found: %v", challengeID, err)

		// Send an error response if mini challenge is not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Mini challenge not found"})
		return
	}

	// Bind the incoming JSON to the mini challenge model
	if err := c.ShouldBindJSON(&challenge); err != nil {
		// Log the error if binding the JSON fails
		log.Printf("Failed to bind request data for mini challenge update: %v", err)

		// Send an error response if JSON binding fails
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated mini challenge to the database
	if err := database.DB.Save(&challenge).Error; err != nil {
		// Log the error if saving the updated mini challenge fails
		log.Printf("Failed to update mini challenge with ID %s: %v", challengeID, err)

		// Send an error response if saving the updated mini challenge fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mini challenge"})
		return
	}

	// Log the successful mini challenge update
	log.Printf("Successfully updated mini challenge with ID: %s", challengeID)

	// Send the updated mini challenge as the response
	c.JSON(http.StatusOK, challenge)
}
