package User

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"fmt"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadProfilePicture(c *gin.Context) {
	userID := c.Param("userId")

	file, err := c.FormFile("profilePic")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// In a real implementation, you would save this to cloud storage
	// For now, we'll just simulate a URL
	profilePicURL := fmt.Sprintf("https://example.com/profiles/user%s_%s", userID, file.Filename)

	// Update user profile picture URL
	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Update("profile_pic_url", profilePicURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile picture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profilePicUrl": profilePicURL})
}
