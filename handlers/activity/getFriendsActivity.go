package activity

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetFriendsActivities(c *gin.Context) {
	var req models.FriendsActivityRequest

	// Log incoming request body
	log.Printf("Received request to get friends' activities: %v", c.Request.Body)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Parsed userIds: %v", req.UserIds)

	var userIDs []uint
	for _, idStr := range req.UserIds {
		if id, err := strconv.Atoi(idStr); err == nil {
			userIDs = append(userIDs, uint(id))
		} else {
			log.Printf("Invalid userId string '%s', skipping.", idStr)
		}
	}

	log.Printf("Final userIDs: %v", userIDs)

	var activities []models.Activity
	// Querying activities from database
	log.Printf("Querying activities for user IDs: %v", userIDs)
	if err := database.DB.Where("user_id IN ?", userIDs).Order("created_at DESC").Limit(20).Find(&activities).Error; err != nil {
		log.Printf("Error fetching activities from DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	log.Printf("Successfully fetched %d activities", len(activities))

	c.JSON(http.StatusOK, activities)
}
