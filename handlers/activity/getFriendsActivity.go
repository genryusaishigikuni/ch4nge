package activity

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetFriendsActivities(c *gin.Context) {
	var req models.FriendsActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userIDs []uint
	for _, idStr := range req.UserIds {
		if id, err := strconv.Atoi(idStr); err == nil {
			userIDs = append(userIDs, uint(id))
		}
	}

	var activities []models.Activity
	if err := db.DB.Where("user_id IN ?", userIDs).Order("created_at DESC").Limit(20).Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}
