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
		Preload("User").       // Load user information
		Preload("PostLikes").  // Load likes
		Preload("PostShares"). // Load shares
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

	// Calculate engagement stats for each post
	type PostWithStats struct {
		models.Post
		LikesCount  int `json:"likes_count"`
		SharesCount int `json:"shares_count"`
	}

	var postsWithStats []PostWithStats
	for _, post := range posts {
		postStats := PostWithStats{
			Post:        post,
			LikesCount:  post.LikeNumber,
			SharesCount: post.SharesNumber,
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
