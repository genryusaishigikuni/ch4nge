package admin

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateWeeklyChallenge(c *gin.Context) {
	var challenge models.WeeklyChallenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&challenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create weekly challenge"})
		return
	}

	c.JSON(http.StatusCreated, challenge)
}
