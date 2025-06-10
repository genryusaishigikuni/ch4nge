package activity

import (
	"log"
	"net/http"
	"strconv"

	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func GetUserActivities(c *gin.Context) {
	// Log the incoming request
	userIdStr := c.Param("userId")
	log.Printf("Received request to get activities for userId: %s", userIdStr)

	// Parsing userId from path parameter
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Printf("Error parsing userId '%s': %v", userIdStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	log.Printf("Parsed userId: %d", userId)

	// Fetch limit and offset query parameters with default values
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Printf("Error parsing limit '%s': %v", limit, err)
		limitInt = 20 // fallback to default limit
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		log.Printf("Error parsing offset '%s': %v", offset, err)
		offsetInt = 0 // fallback to default offset
	}

	log.Printf("Query parameters: limit=%d, offset=%d", limitInt, offsetInt)

	// Query activities from database
	var activities []models.Activity
	log.Printf("Fetching activities for userId=%d with limit=%d and offset=%d", userId, limitInt, offsetInt)
	if err := database.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(limitInt).
		Offset(offsetInt).
		Find(&activities).Error; err != nil {
		log.Printf("Error fetching activities from database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user activities"})
		return
	}

	log.Printf("Successfully fetched %d activities for userId=%d", len(activities), userId)

	// Count the total activities for pagination
	var totalCount int64
	if err := database.DB.Model(&models.Activity{}).Where("user_id = ?", userId).Count(&totalCount).Error; err != nil {
		log.Printf("Error counting total activities for userId=%d: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total activities"})
		return
	}
	log.Printf("Total activities count for userId=%d: %d", userId, totalCount)

	// Construct the response
	response := gin.H{
		"activities": activities,
		"total":      totalCount,
		"limit":      limitInt,
		"offset":     offsetInt,
		"has_more":   int64(offsetInt+limitInt) < totalCount,
	}

	// Send the response
	log.Printf("Sending response: %v", response)
	c.JSON(http.StatusOK, response)
}
