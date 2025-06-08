package post

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"fmt"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func UploadPost(c *gin.Context) {
	var req models.PostRequest

	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	post := models.Post{
		UserID: uint(userID),
		Title:  req.Title,
	}

	if file, err := c.FormFile("image"); err == nil {
		// In a real implementation, save to cloud storage
		imageURL := fmt.Sprintf("https://example.com/posts/%d_%s", userID, file.Filename)
		post.ImageURL = imageURL
	}

	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}
