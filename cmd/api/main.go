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

	// Run database migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize database
	db := database.InitDB()

	// Get port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	// Apply CORS middleware first
	router.Use(middleware.CORSMiddleware())

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
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", api.register)
			auth.POST("/login", api.login)
			auth.POST("/oldLogin", api.oldLogin)
			auth.POST("/google-auth", api.googleAuth)
			auth.POST("/check-mobile", api.checkMobileNumber)
		}

		// Protected routes
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
