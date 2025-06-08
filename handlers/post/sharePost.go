package post

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/genryusaishigikuni/ch4nge/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func SharePost(c *gin.Context) {
	postID := c.Param("postId")

	var req models.SharePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	share := models.PostShare{
		PostID:   utils.ParseUint(postID),
		UserID:   uint(userID),
		SharedAt: time.Now(),
	}
	if err := db.DB.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share post"})
		return
	}

	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	post.SharedBy = append(post.SharedBy, uint(userID))
	if err := db.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post share list"})
		return
	}

	c.JSON(http.StatusOK, post)
}
