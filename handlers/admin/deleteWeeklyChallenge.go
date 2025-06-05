package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteWeeklyChallenge(c *gin.Context) {
	challengeID := c.Param("id")

	if err := db.DB.Delete(&models.WeeklyChallenge{}, challengeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete weekly challenge"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Weekly challenge deleted successfully"})
}
