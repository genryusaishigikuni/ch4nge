package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateMiniChallenge(c *gin.Context) {
	challengeID := c.Param("id")

	var challenge models.MiniChallenge
	if err := db.DB.First(&challenge, challengeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mini challenge not found"})
		return
	}

	if err := c.ShouldBindJSON(&challenge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&challenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mini challenge"})
		return
	}

	c.JSON(http.StatusOK, challenge)
}
