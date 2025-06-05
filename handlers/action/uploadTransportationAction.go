package action

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func UploadTransportationAction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.TransportationActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate points based on eco-friendliness
	points := 5 // Default points
	if metadata, ok := req.Metadata["isEcoFriendly"].(bool); ok && metadata {
		points = 15
	}

	action := models.Action{
		UserID:     userID.(uint),
		ActionType: req.ActionType,
		Payload:    req.Payload,
		Metadata:   req.Metadata,
		Points:     points,
	}

	if err := db.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload transportation action"})
		return
	}

	// Update user points
	db.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", action.Points))

	// Create activity
	activity := models.Activity{
		UserID: userID.(uint),
		Title:  "Used transportation",
		Value:  points,
	}
	db.DB.Create(&activity)

	c.JSON(http.StatusCreated, gin.H{"message": "Transportation action uploaded successfully."})
}
