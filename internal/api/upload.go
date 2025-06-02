package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUploadRoutes(r *gin.RouterGroup) {
	upload := r.Group("/upload")
	{
		upload.POST("", uploadImage)
	}
}

func uploadImage(c *gin.Context) {
	// Get the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	if !isValidImageType(ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only images are allowed"})
		return
	}

	// Validate file size (5MB limit)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum size is 5MB"})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Save the file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Return the URL
	url := fmt.Sprintf("/uploads/%s", filename)
	log.Printf("Avatar uploaded: %s", url)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func isValidImageType(ext string) bool {
	validTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return validTypes[ext]
}
