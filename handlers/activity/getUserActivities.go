package activity

import (
	"net/http"
	"strconv"

	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func GetUserActivities(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	var activities []models.Activity
	if err := database.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(limitInt).
		Offset(offsetInt).
		Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user activities"})
		return
	}

	var totalCount int64
	database.DB.Model(&models.Activity{}).Where("user_id = ?", userId).Count(&totalCount)

	response := gin.H{
		"activities": activities,
		"total":      totalCount,
		"limit":      limitInt,
		"offset":     offsetInt,
		"has_more":   int64(offsetInt+limitInt) < totalCount,
	}

	c.JSON(http.StatusOK, response)
}
