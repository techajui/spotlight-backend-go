package middleware

import (
	"fmt"
	"time"

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
