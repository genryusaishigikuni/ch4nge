package post

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"fmt"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func UploadPost(c *gin.Context) {
	var req models.PostRequest

	// Handle different content types
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := c.ShouldBind(&req); err != nil {
			log.Printf("ERROR: Failed to bind multipart form data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("ERROR: Failed to bind JSON data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Validate required fields
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		log.Printf("ERROR: Invalid user ID: %s", req.UserID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	log.Printf("DEBUG: Creating post for user ID: %d", userID)

	// Create post object
	post := models.Post{
		UserID: uint(userID),
		Title:  req.Title,
	}

	// Handle optional description
	if req.Description != "" {
		post.Description = req.Description
	}

	// Handle image upload if present
	if file, err := c.FormFile("image"); err == nil {
		log.Printf("DEBUG: Image file received: %s, Size: %d bytes", file.Filename, file.Size)

		// Validate image
		if err := validatePostImage(file); err != nil {
			log.Printf("ERROR: Image validation failed: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Upload image and get URL
		imageURL, err := savePostImage(file, uint(userID))
		if err != nil {
			log.Printf("ERROR: Failed to save post image: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		post.ImageURL = imageURL
		log.Printf("DEBUG: Post image saved with URL: %s", imageURL)
	} else {
		log.Printf("DEBUG: No image file provided or error getting file: %v", err)
	}

	// Save post to database
	if err := db.DB.Create(&post).Error; err != nil {
		log.Printf("ERROR: Failed to create post in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	log.Printf("DEBUG: Post created successfully with ID: %d", post.ID)

	// Return success response with post details
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Post created successfully",
		"post_id":    post.ID,
		"title":      post.Title,
		"image_url":  post.ImageURL,
		"created_at": post.CreatedAt,
	})
}

func validatePostImage(file *multipart.FileHeader) error {
	// Check file size (max 10MB for posts)
	const maxSize = 10 * 1024 * 1024 // 10MB
	if file.Size > maxSize {
		return fmt.Errorf("image too large (max 10MB), got %d bytes", file.Size)
	}

	// Check for empty file
	if file.Size == 0 {
		return fmt.Errorf("image file is empty")
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
		return fmt.Errorf("invalid image type '%s'. Only JPG, JPEG, PNG, GIF, and WEBP are allowed", ext)
	}

	// Check filename
	if file.Filename == "" {
		return fmt.Errorf("filename is empty")
	}

	return nil
}

func savePostImage(file *multipart.FileHeader, userID uint) (string, error) {
	// Get current working directory
	cwd, _ := os.Getwd()
	log.Printf("DEBUG: Current working directory: %s", cwd)

	// Define upload directory path for posts
	uploadDir := filepath.Join(cwd, "uploads", "posts")
	log.Printf("DEBUG: Post upload directory path: %s", uploadDir)

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("ERROR: Failed to create post upload directory: %v", err)
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}
	log.Printf("DEBUG: Post directory created/verified successfully")

	// Check if directory is writable
	testFile := filepath.Join(uploadDir, ".test_write")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		log.Printf("ERROR: Post upload directory not writable: %v", err)
		return "", fmt.Errorf("upload directory not writable: %v", err)
	}
	err := os.Remove(testFile)
	if err != nil {
		return "", err
	} // Clean up test file
	log.Printf("DEBUG: Post directory write test passed")

	// Generate unique filename
	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("post_%d_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)
	log.Printf("DEBUG: Generated post file path: %s", filePath)

	// Save the uploaded file
	log.Printf("DEBUG: Saving post image to: %s", filePath)
	if err := saveUploadedFile(file, filePath); err != nil {
		log.Printf("ERROR: Failed to save post image: %v", err)
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	log.Printf("DEBUG: Post image saved successfully")

	// Verify file was saved correctly
	if info, err := os.Stat(filePath); err != nil {
		log.Printf("ERROR: Post image file verification failed: %v", err)
		return "", fmt.Errorf("file save verification failed")
	} else {
		log.Printf("DEBUG: Post image file verified - Size: %d bytes", info.Size())
	}

	// Create URL path for the uploaded file
	imageURL := fmt.Sprintf("/uploads/posts/%s", filename)
	log.Printf("DEBUG: Generated post image URL: %s", imageURL)

	return imageURL, nil
}

// Helper function to save uploaded file (similar to Gin's SaveUploadedFile but with error handling)
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Printf("ERROR: Failed to close file: %v", err)
		}
	}(src)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Printf("ERROR: Failed to close file: %v", err)
		}
	}(out)

	_, err = io.Copy(out, src)
	return err
}
