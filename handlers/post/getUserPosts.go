package post

import (
	"net/http"
	"strconv"

	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func GetUserPosts(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get query parameters for pagination
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	var posts []models.Post
	if err := database.DB.Where("user_id = ?", userId).
		Preload("User"). // Load user information
		Order("created_at DESC").
		Limit(limitInt).
		Offset(offsetInt).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user posts"})
		return
	}

	// Get total count for pagination
	var totalCount int64
	database.DB.Model(&models.Post{}).Where("user_id = ?", userId).Count(&totalCount)

	// Enhanced post response with engagement stats and user interaction info
	type PostWithStats struct {
		models.Post
		LikesCount    int  `json:"likes_count"`
		SharesCount   int  `json:"shares_count"`
		IsLikedByUser bool `json:"is_liked_by_user,omitempty"` // Optional: if you want to check for requesting user
	}

	var postsWithStats []PostWithStats
	for _, post := range posts {
		postStats := PostWithStats{
			Post:        post,
			LikesCount:  len(post.LikedBy),
			SharesCount: len(post.SharedBy),
		}

		// Get requesting user ID from auth context (optional)
		var requestingUserID uint
		if userIDInterface, exists := c.Get("userID"); exists {
			if id, ok := userIDInterface.(uint); ok {
				requestingUserID = id
			}
		}

		// Check if the requesting user has liked this post
		if requestingUserID > 0 {
			postStats.IsLikedByUser = post.HasLiked(requestingUserID)
		}

		postsWithStats = append(postsWithStats, postStats)
	}

	response := gin.H{
		"posts":    postsWithStats,
		"total":    totalCount,
		"limit":    limitInt,
		"offset":   offsetInt,
		"has_more": int64(offsetInt+limitInt) < totalCount,
	}

	c.JSON(http.StatusOK, response)
}
