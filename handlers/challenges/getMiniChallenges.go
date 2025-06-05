package challenges

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMiniChallenges(c *gin.Context) {
	userID := c.Param("userId")

	var userChallenges []models.UserMiniChallenge
	if err := db.DB.Preload("MiniChallenge").Where("user_id = ?", userID).Find(&userChallenges).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mini challenges"})
		return
	}

	var responses []models.MiniChallengeResponse
	for _, uc := range userChallenges {
		responses = append(responses, models.MiniChallengeResponse{
			MiniChallengeID: uc.MiniChallenge.ID,
			UserID:          uc.UserID,
			Title:           uc.MiniChallenge.Title,
			Subtitle:        uc.MiniChallenge.Subtitle,
			IsAchieved:      uc.IsAchieved,
			Points:          uc.MiniChallenge.Points,
		})
	}

	c.JSON(http.StatusOK, responses)
}
