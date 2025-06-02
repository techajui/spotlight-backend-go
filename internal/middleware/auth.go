package middleware

import (
	"fmt"
	"net/http"
	"os"
	"spotlight-backend-go/internal/database"
	"spotlight-backend-go/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("Auth header: %s\n", authHeader) // Debug log

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("Invalid auth format: %v\n", parts) // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		fmt.Printf("Token string: %s\n", tokenString) // Debug log

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				panic("JWT_SECRET environment variable is not set")
			}
			fmt.Printf("JWT Secret: %s\n", jwtSecret) // Debug log
			return []byte(jwtSecret), nil
		})

		if err != nil {
			fmt.Printf("Token parse error: %v\n", err) // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			fmt.Printf("Token is invalid\n") // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Printf("Invalid claims type\n") // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		fmt.Printf("Claims: %v\n", claims) // Debug log

		userID, ok := claims["user_id"].(string)
		if !ok {
			fmt.Printf("Invalid user_id in claims\n") // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
			c.Abort()
			return
		}

		// Trim any whitespace from the user ID
		userID = strings.TrimSpace(userID)

		var user models.User
		if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
			fmt.Printf("User not found: %v\n", err) // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}
