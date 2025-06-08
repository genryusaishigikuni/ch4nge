package User

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"fmt"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UploadProfilePicture(c *gin.Context) {
	userID := c.Param("userId")
	log.Printf("DEBUG: Starting upload for user ID: %s", userID)

	file, err := c.FormFile("profilePic")
	if err != nil {
		log.Printf("ERROR: No file uploaded: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	log.Printf("DEBUG: File received: %s, Size: %d bytes", file.Filename, file.Size)

	// Validate file
	if err := validateImage(file); err != nil {
		log.Printf("ERROR: File validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("DEBUG: File validation passed")

	// Check current working directory
	cwd, _ := os.Getwd()
	log.Printf("DEBUG: Current working directory: %s", cwd)

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads/profiles"
	log.Printf("DEBUG: Attempting to create directory: %s", uploadDir)

	if err := os.MkdirAll(uploadDir, 0777); err != nil {
		log.Printf("ERROR: Failed to create upload directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create upload directory: %v", err)})
		return
	}
	log.Printf("DEBUG: Directory created successfully")

	// Check if directory exists and is writable
	if info, err := os.Stat(uploadDir); err != nil {
		log.Printf("ERROR: Cannot stat upload directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload directory not accessible"})
		return
	} else {
		log.Printf("DEBUG: Directory exists, IsDir: %v, Mode: %v", info.IsDir(), info.Mode())
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("user_%s_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)
	log.Printf("DEBUG: Generated file path: %s", filePath)

	// Save file
	log.Printf("DEBUG: Attempting to save file to: %s", filePath)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("ERROR: Failed to save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
		return
	}
	log.Printf("DEBUG: File saved successfully")

	// Verify file was saved
	if _, err := os.Stat(filePath); err != nil {
		log.Printf("ERROR: File not found after save: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File save verification failed"})
		return
	}
	log.Printf("DEBUG: File save verified")

	// Create URL that will be accessible from outside
	profilePicURL := fmt.Sprintf("/uploads/profiles/%s", filename)
	log.Printf("DEBUG: Generated profile pic URL: %s", profilePicURL)

	// Update user profile picture URL in database
	log.Printf("DEBUG: Updating database for user ID: %s", userID)
	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Update("profile_pic_url", profilePicURL).Error; err != nil {
		log.Printf("ERROR: Failed to update profile picture in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update profile picture: %v", err)})
		return
	}
	log.Printf("DEBUG: Database updated successfully")

	log.Printf("DEBUG: Upload completed successfully for user %s", userID)
	c.JSON(http.StatusOK, gin.H{"profilePicUrl": profilePicURL})
}

func validateImage(file *multipart.FileHeader) error {
	// Check file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return fmt.Errorf("file too large (max 5MB)")
	}

	// Check file extension (case-insensitive)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("invalid file type. Only JPG, PNG, and GIF allowed")
	}

	return nil
}
