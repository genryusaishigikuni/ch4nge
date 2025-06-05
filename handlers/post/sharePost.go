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

	// Create share
	share := models.PostShare{
		PostID: utils.ParseUint(postID),
		UserID: uint(userID),
	}
	if err := db.DB.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share post"})
		return
	}

	// Update share count
	db.DB.Model(&models.Post{}).Where("id = ?", postID).Update("shares_number", gorm.Expr("shares_number + 1"))

	// Return updated post
	var post models.Post
	db.DB.First(&post, postID)
	c.JSON(http.StatusOK, post)
}
