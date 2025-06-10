package challenges

import (
	"log"
	"net/http"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func GetMiniChallenges(c *gin.Context) {
	userID := c.Param("userId")
	log.Printf("Fetching mini challenges for user ID: %s", userID) // Log user ID for which mini challenges are being fetched

	// Fetch the mini challenges for the user
	var userChallenges []models.UserMiniChallenge
	if err := db.DB.Preload("MiniChallenge").Where("user_id = ?", userID).Find(&userChallenges).Error; err != nil {
		log.Printf("Error fetching mini challenges for user ID: %s, error: %v", userID, err) // Log any errors encountered
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mini challenges"})
		return
	}

	// Log the count of challenges fetched
	log.Printf("Fetched %d mini challenges for user ID: %s", len(userChallenges), userID)

	var responses []models.MiniChallengeResponse
	for _, uc := range userChallenges {
		// Log the individual challenge being processed
		log.Printf("Processing mini challenge ID: %d, Title: %s", uc.MiniChallenge.ID, uc.MiniChallenge.Title)

		// Append the mini challenge details to the response array
		responses = append(responses, models.MiniChallengeResponse{
			MiniChallengeID: uc.MiniChallenge.ID,
			UserID:          uc.UserID,
			Title:           uc.MiniChallenge.Title,
			Subtitle:        uc.MiniChallenge.Subtitle,
			IsAchieved:      uc.IsAchieved,
			Points:          uc.MiniChallenge.Points,
		})
	}

	// Log the response data before sending it back to the client
	log.Printf("Sending response with %d mini challenges for user ID: %s", len(responses), userID)

	// Send the mini challenges as a JSON response
	c.JSON(http.StatusOK, responses)
}
