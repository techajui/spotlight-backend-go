package middleware

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",          // React dev server
			"http://localhost:5173",          // Vite (optional)
			"http://localhost:8080",          // Backend (if used directly)
			"https://spot.smartrating.in",    // Production frontend
			"https://spotlight-backend-go.onrender.com", // Production API
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length", "Content-Type",
		},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// Accept all localhost origins for local dev flexibility
			return origin == "http://localhost:3000" || 
				origin == "http://localhost:8080" || 
				origin == "https://spot.smartrating.in" ||
				origin == "https://spotlight-backend-go.onrender.com"
		},
		MaxAge: 12 * time.Hour,
	})
}

// RateLimitMiddleware returns a rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
	// Create a rate limiter with a limit of 100 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	store := memory.NewStore()
	limiter := limiter.New(store, rate)

	return func(c *gin.Context) {
		context := c.Request.Context()
		ip := c.ClientIP()
		key := fmt.Sprintf("%s:%s", ip, c.Request.URL.Path)

		// Check if the request is allowed
		if _, err := limiter.Get(context, key); err != nil {
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
