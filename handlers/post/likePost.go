package post

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/genryusaishigikuni/ch4nge/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func LikePost(c *gin.Context) {
	postID := c.Param("postId")

	var req models.LikePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if already liked
	var existingLike models.PostLike
	if err := db.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Post already liked"})
		return
	}

	// Create like
	like := models.PostLike{
		PostID: utils.ParseUint(postID),
		UserID: uint(userID),
	}
	if err := db.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	// Update like count
	db.DB.Model(&models.Post{}).Where("id = ?", postID).Update("like_number", gorm.Expr("like_number + 1"))

	// Return updated post
	var post models.Post
	db.DB.First(&post, postID)
	c.JSON(http.StatusOK, post)
}
