package admin

import (
	"log"
	"net/http"
	"time"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func AssignMiniChallengeToUser(c *gin.Context) {
	// Log the incoming request
	log.Println("Received request to assign mini challenge to user")

	// Define the request structure
	var req struct {
		UserID          uint `json:"userId" binding:"required"`
		MiniChallengeID uint `json:"miniChallengeId" binding:"required"`
	}

	// Bind the incoming JSON request to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Request successfully bound. UserID: %d, MiniChallengeID: %d", req.UserID, req.MiniChallengeID)

	// Create the UserMiniChallenge record
	userChallenge := models.UserMiniChallenge{
		UserID:          req.UserID,
		MiniChallengeID: req.MiniChallengeID,
		AssignedAt:      time.Now(),
	}

	// Save to the database
	log.Printf("Attempting to assign mini challenge %d to user %d", req.MiniChallengeID, req.UserID)
	if err := db.DB.Create(&userChallenge).Error; err != nil {
		log.Printf("Error assigning mini challenge: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign mini challenge"})
		return
	}

	// Success response
	log.Printf("Mini challenge %d successfully assigned to user %d", req.MiniChallengeID, req.UserID)
	c.JSON(http.StatusCreated, gin.H{"message": "Mini challenge assigned successfully"})
}
