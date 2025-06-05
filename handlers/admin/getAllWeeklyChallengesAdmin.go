package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllWeeklyChallengesAdmin(c *gin.Context) {
	var challenges []models.WeeklyChallenge
	if err := db.DB.Find(&challenges).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weekly challenges"})
		return
	}

	c.JSON(http.StatusOK, challenges)
}
