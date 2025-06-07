package action

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func UploadGreenAction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.GreenActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action := models.Action{
		UserID:     userID.(uint),
		ActionType: req.ActionType,
		Payload:    req.Payload,
		Metadata:   req.Metadata,
		Points:     10, // Default points for green actions
	}

	if err := db.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload green action"})
		return
	}

	// Update user points
	db.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", action.Points))

	// Update user location if provided in action payload
	if locationArray, ok := req.Payload["location"].([]interface{}); ok && len(locationArray) == 2 {
		if latitude, ok := locationArray[0].(float64); ok {
			if longitude, ok := locationArray[1].(float64); ok {
				db.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
					"latitude":  latitude,
					"longitude": longitude,
				})
			}
		}
	}

	// Get action title from payload
	actionTitle := "Performed a green action"
	if option, ok := req.Payload["option"].(string); ok && option != "" {
		actionTitle = formatGreenActionTitle(option)
	}

	// Create activity
	activity := models.Activity{
		UserID: userID.(uint),
		Title:  actionTitle,
		Value:  action.Points,
	}
	db.DB.Create(&activity)

	c.JSON(http.StatusCreated, gin.H{"message": "Green action uploaded successfully."})
}
