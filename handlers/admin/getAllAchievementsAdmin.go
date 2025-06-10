package admin

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetAllAchievementsAdmin(c *gin.Context) {
	// Log the request to fetch all achievements
	log.Println("Received request to fetch all achievements")

	// Retrieve all achievements from the database
	var achievements []models.Achievement
	if err := database.DB.Find(&achievements).Error; err != nil {
		// Log the error if fetching the achievements fails
		log.Printf("Error fetching achievements: %v", err)

		// Send an error response back to the client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
		return
	}

	// Log successful achievement retrieval
	log.Printf("Successfully fetched %d achievements", len(achievements))

	// Send the retrieved achievements back to the client
	c.JSON(http.StatusOK, achievements)
}
