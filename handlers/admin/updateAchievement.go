package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateAchievement(c *gin.Context) {
	achievementID := c.Param("id")

	var achievement models.Achievement
	if err := db.DB.First(&achievement, achievementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}

	if err := c.ShouldBindJSON(&achievement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&achievement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievement"})
		return
	}

	c.JSON(http.StatusOK, achievement)
}
