package challenges

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWeeklyChallenge(c *gin.Context) {
	userID := c.Param("userId")

	var userChallenge models.UserWeeklyChallenge
	if err := db.DB.Preload("WeeklyChallenge").Where("user_id = ?", userID).First(&userChallenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No weekly challenge found"})
		return
	}

	response := models.WeeklyChallengeResponse{
		WeeklyChallengeID: userChallenge.WeeklyChallenge.ID,
		UserID:            userChallenge.UserID,
		Title:             userChallenge.WeeklyChallenge.Title,
		Subtitle:          userChallenge.WeeklyChallenge.Subtitle,
		CurrentValue:      userChallenge.CurrentValue,
		TotalValue:        userChallenge.WeeklyChallenge.TargetValue,
		Points:            userChallenge.WeeklyChallenge.Points,
	}

	c.JSON(http.StatusOK, response)
}
