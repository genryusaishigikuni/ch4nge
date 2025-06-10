package admin

import (
	"log"
	"net/http"
	"time"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func AssignWeeklyChallengeToUser(c *gin.Context) {
	// Log the incoming request
	log.Println("Received request to assign weekly challenge to user")

	// Define the request structure
	var req models.UserWeeklyChallengeRequest

	// Bind the incoming JSON request to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Request successfully bound. UserID: %d, WeeklyChallengeID: %d, CurrentValue: %f", req.UserID, req.WeeklyChallengeID, req.CurrentValue)

	// Create the UserWeeklyChallenge record
	userChallenge := models.UserWeeklyChallenge{
		UserID:            req.UserID,
		WeeklyChallengeID: req.WeeklyChallengeID,
		CurrentValue:      req.CurrentValue,
		AssignedAt:        time.Now(),
	}

	// Save to the database
	log.Printf("Attempting to assign weekly challenge %d to user %d", req.WeeklyChallengeID, req.UserID)
	if err := db.DB.Create(&userChallenge).Error; err != nil {
		log.Printf("Error assigning weekly challenge: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign weekly challenge"})
		return
	}

	// Success response
	log.Printf("Weekly challenge %d successfully assigned to user %d", req.WeeklyChallengeID, req.UserID)
	c.JSON(http.StatusCreated, gin.H{"message": "Weekly challenge assigned successfully"})
}
