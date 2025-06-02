package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"spotlight-backend-go/internal/database"
	"spotlight-backend-go/internal/models"
	"spotlight-backend-go/internal/schemas"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

func RegisterAuthRoutes(r *gin.RouterGroup) {
	r.POST("/register", register)
	r.POST("/login", login)
}

func register(c *gin.Context) {
	var req schemas.UserCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// Convert MediaGallery to JSON
	mediaGalleryJSON, err := json.Marshal(req.MediaGallery)
	if err != nil {
		mediaGalleryJSON = []byte("[]")
	}

	// Generate username from email (take part before @ and add random number)
	username := strings.Split(req.Email, "@")[0]
	username = strings.ToLower(username)
	// Add random number to ensure uniqueness
	username = fmt.Sprintf("%s%d", username, time.Now().UnixNano()%10000)

	user := models.User{
		ID:           generateUUID(),
		Name:         req.Name,
		Username:     username,
		Email:        req.Email,
		Password:     string(hashedPassword),
		Role:         models.RoleFan, // Default to fan, or allow from req if needed
		MediaGallery: datatypes.JSON(mediaGalleryJSON),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Transform user to match frontend expectations
	response := gin.H{
		"user": gin.H{
			"id":              user.ID,
			"name":            user.Name,
			"username":        user.Username,
			"email":           user.Email,
			"role":            user.Role,
			"avatar_url":      user.AvatarURL,
			"bio":             user.Bio,
			"mediaGallery":    user.MediaGallery,
			"walletBalance":   user.WalletBalance,
			"followerCount":   user.FollowerCount,
			"instagramHandle": user.InstagramHandle,
			"verified":        user.Verified,
		},
	}

	c.JSON(http.StatusCreated, response)
}

func login(c *gin.Context) {
	var req schemas.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token (implement your own token logic)
	token := generateToken(user.ID, user.Role)

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// Dummy UUID generator (replace with a real one)
func generateUUID() string {
	return fmt.Sprintf("user-%d", time.Now().UnixNano())
}

// Generate a proper JWT token
func generateToken(userID string, role models.UserRole) string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is not set")
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic("Failed to generate token: " + err.Error())
	}

	return tokenString
}
