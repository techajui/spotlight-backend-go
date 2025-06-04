package main

import (
	"log"
	"os"
	"spotlight-backend-go/internal/api"
	"spotlight-backend-go/internal/database"
	"spotlight-backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db := database.InitDB()

	// Get port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	// ✅ Apply CORS early — critical for OPTIONS requests
	router.Use(middleware.CORSMiddleware())

	// Optional: manually handle all OPTIONS requests (needed if not using gin-contrib/cors preflight)
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})

	// Security & Rate limiting — skip in local dev
	if os.Getenv("LOCAL_DEV") != "true" {
		router.Use(middleware.SecurityMiddleware())
		router.Use(middleware.RateLimitMiddleware())
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve uploaded files
	router.Static("/uploads", "./uploads")

	// API v1
	v1 := router.Group("/api/v1")
	{
		api.RegisterAuthRoutes(v1)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			api.RegisterUserRoutes(protected)
			api.RegisterEventRoutes(protected, db)
			api.RegisterChatRoutes(protected, db)
			api.RegisterUploadRoutes(protected)
			api.RegisterApplicationRoutes(protected, db)
		}
	}

	// Start server
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
