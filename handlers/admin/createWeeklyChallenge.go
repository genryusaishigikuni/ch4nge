package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreateWeeklyChallenge(c *gin.Context) {
	var req models.WeeklyChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	challenge := models.WeeklyChallenge{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Points:      req.Points,
		TargetValue: req.TargetValue,
		IsActive:    req.IsActive,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := db.DB.Create(&challenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create weekly challenge"})
		return
	}

	// Assign to all existing users
	if err := assignWeeklyChallengeToAllUsers(challenge.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Weekly challenge created but failed to assign to all users"})
		return
	}

	c.JSON(http.StatusCreated, challenge)
}

func assignWeeklyChallengeToAllUsers(challengeID uint) error {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		userChallenge := models.UserWeeklyChallenge{
			UserID:            user.ID,
			WeeklyChallengeID: challengeID,
			CurrentValue:      0.0,
			AssignedAt:        time.Now(),
		}
		if err := db.DB.FirstOrCreate(&userChallenge, models.UserWeeklyChallenge{
			UserID:            user.ID,
			WeeklyChallengeID: challengeID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
