package User

import (
	"net/http"
	"strconv"

	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func GetUserActions(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get query parameters for pagination and filtering
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")
	actionType := c.Query("type") // Optional filter by action type

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	var actions []models.Action
	query := database.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(limitInt).
		Offset(offsetInt)

	// Apply type filter if provided
	if actionType != "" {
		query = query.Where("action_type = ?", actionType)
	}

	if err := query.Find(&actions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user actions"})
		return
	}

	// Get total count for pagination
	var totalCount int64
	countQuery := database.DB.Model(&models.Action{}).Where("user_id = ?", userId)
	if actionType != "" {
		countQuery = countQuery.Where("action_type = ?", actionType)
	}
	countQuery.Count(&totalCount)

	response := gin.H{
		"actions":  actions,
		"total":    totalCount,
		"limit":    limitInt,
		"offset":   offsetInt,
		"has_more": int64(offsetInt+limitInt) < totalCount,
	}

	c.JSON(http.StatusOK, response)
}
