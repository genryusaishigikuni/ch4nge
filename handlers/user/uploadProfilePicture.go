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

	// Get the uploaded file
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

	// Get current working directory
	cwd, _ := os.Getwd()
	log.Printf("DEBUG: Current working directory: %s", cwd)

	// Define upload directory path
	uploadDir := filepath.Join(cwd, "uploads", "profiles")
	log.Printf("DEBUG: Upload directory path: %s", uploadDir)

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("ERROR: Failed to create upload directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create upload directory",
			"details": err.Error(),
		})
		return
	}
	log.Printf("DEBUG: Directory created/verified successfully")

	// Check if directory is writable
	testFile := filepath.Join(uploadDir, ".test_write")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		log.Printf("ERROR: Upload directory not writable: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Upload directory not writable",
			"details": err.Error(),
		})
		return
	}
	os.Remove(testFile) // Clean up test file
	log.Printf("DEBUG: Directory write test passed")

	// Generate unique filename
	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("user_%s_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)
	log.Printf("DEBUG: Generated file path: %s", filePath)

	// Delete old profile picture if exists
	var user models.User
	if err := db.DB.First(&user, userID).Error; err == nil && user.ProfilePicURL != "" {
		// Extract filename from URL and delete old file
		oldFilename := filepath.Base(user.ProfilePicURL)
		oldFilePath := filepath.Join(uploadDir, oldFilename)
		if _, err := os.Stat(oldFilePath); err == nil {
			if err := os.Remove(oldFilePath); err != nil {
				log.Printf("WARNING: Failed to delete old profile picture: %v", err)
			} else {
				log.Printf("DEBUG: Deleted old profile picture: %s", oldFilePath)
			}
		}
	}

	// Save the uploaded file
	log.Printf("DEBUG: Saving file to: %s", filePath)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("ERROR: Failed to save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save file",
			"details": err.Error(),
		})
		return
	}
	log.Printf("DEBUG: File saved successfully")

	// Verify file was saved correctly
	if info, err := os.Stat(filePath); err != nil {
		log.Printf("ERROR: File verification failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File save verification failed"})
		return
	} else {
		log.Printf("DEBUG: File verified - Size: %d bytes", info.Size())
	}

	// Create URL path for the uploaded file
	profilePicURL := fmt.Sprintf("/uploads/profiles/%s", filename)
	log.Printf("DEBUG: Generated profile pic URL: %s", profilePicURL)

	// Update user profile picture URL in database
	log.Printf("DEBUG: Updating database for user ID: %s", userID)
	result := db.DB.Model(&models.User{}).Where("id = ?", userID).Update("profile_pic_url", profilePicURL)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update profile picture in database: %v", result.Error)
		// Clean up uploaded file on database error
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile picture in database",
			"details": result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("ERROR: No user found with ID: %s", userID)
		// Clean up uploaded file
		os.Remove(filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	log.Printf("DEBUG: Database updated successfully, rows affected: %d", result.RowsAffected)
	log.Printf("DEBUG: Upload completed successfully for user %s", userID)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":       "Profile picture uploaded successfully",
		"profilePicUrl": profilePicURL,
		"filename":      filename,
		"fileSize":      file.Size,
	})
}

func validateImage(file *multipart.FileHeader) error {
	// Check file size (max 5MB)
	const maxSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxSize {
		return fmt.Errorf("file too large (max 5MB), got %d bytes", file.Size)
	}

	// Check for empty file
	if file.Size == 0 {
		return fmt.Errorf("file is empty")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("invalid file type '%s'. Only JPG, JPEG, PNG, GIF, and WEBP are allowed", ext)
	}

	// Check filename
	if file.Filename == "" {
		return fmt.Errorf("filename is empty")
	}

	return nil
}
