package post

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/genryusaishigikuni/ch4nge/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LikePost(c *gin.Context) {
	postID := c.Param("postId")

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in token"})
		return
	}

	var post models.Post
	if err := db.DB.First(&post, utils.ParseUint(postID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	isLiked := post.ToggleLike(userID)

	if err := db.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	db.DB.Preload("User").First(&post, post.ID)

	response := gin.H{
		"post":    post,
		"isLiked": isLiked,
		"message": getMessage(isLiked),
	}

	c.JSON(http.StatusOK, response)
}

func getMessage(isLiked bool) string {
	if isLiked {
		return "Post liked successfully"
	}
	return "Post unliked successfully"
}
