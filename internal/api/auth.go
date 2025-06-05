package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/oldLogin", oldLogin)
		auth.POST("/google-auth", googleAuth)
		auth.POST("/check-mobile", checkMobileNumber)
		auth.GET("/test/users", getAllUsers)
	}
}

func register(c *gin.Context) {
	var req schemas.UserCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Registration request binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Attempting registration for email: %s", req.Email)

	// Check if user already exists
	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		log.Printf("User already exists with email: %s", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	log.Printf("Generating password hash for email: %s", req.Email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Convert MediaGallery to JSON
	mediaGalleryJSON, err := json.Marshal(req.MediaGallery)
	if err != nil {
		mediaGalleryJSON = []byte("[]")
	}

	// Convert Interests to JSON
	interestsJSON, err := json.Marshal(req.Interests)
	if err != nil {
		interestsJSON = []byte("[]")
	}

	// Generate username from email (take part before @ and add random number)
	username := strings.Split(req.Email, "@")[0]
	username = strings.ToLower(username)
	// Add random number to ensure uniqueness
	username = fmt.Sprintf("%s%d", username, time.Now().UnixNano()%10000)

	user := models.User{
		ID:             generateUUID(),
		Name:           req.Name,
		Username:       username,
		Email:          req.Email,
		Password:       string(hashedPassword),
		Role:           req.Role,
		Bio:            req.Bio,
		Gender:         req.Gender,
		Age:            req.Age,
		Work:           req.Work,
		Education:      req.Education,
		Interests:      datatypes.JSON(interestsJSON),
		AvatarURL:      req.AvatarURL,
		MediaGallery:   datatypes.JSON(mediaGalleryJSON),
		InstagramHandle: req.InstagramHandle,
		FollowerCount:  req.FollowerCount,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	log.Printf("Creating user in database for email: %s", req.Email)
	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("Database error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	log.Printf("User created successfully for email: %s", req.Email)
	// Transform user to match frontend expectations
	response := gin.H{
		"user": gin.H{
			"id":              user.ID,
			"name":            user.Name,
			"username":        user.Username,
			"email":           user.Email,
			"role":            user.Role,
			"bio":             user.Bio,
			"gender":          user.Gender,
			"age":             user.Age,
			"work":            user.Work,
			"education":       user.Education,
			"interests":       user.Interests,
			"avatar_url":      user.AvatarURL,
			"media_gallery":   user.MediaGallery,
			"instagram_handle": user.InstagramHandle,
			"follower_count":  user.FollowerCount,
			"verified":        user.Verified,
		},
	}

	c.JSON(http.StatusCreated, response)
}

func login(c *gin.Context) {
	var req struct {
		MobileNumber string `json:"mobile_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("mobile_number = ?", req.MobileNumber).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "Please complete your registration",
				"action":  "signup",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate token
	token := generateToken(user.ID, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
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
	})
}

func googleAuth(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Create OAuth2 config
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Create OAuth2 client
	client := config.Client(context.Background(), nil)

	// Create OAuth2 service
	oauth2Service, err := oauth2api.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OAuth2 service"})
		return
	}

	// Try to verify as Google OAuth token first
	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(req.IDToken).Do()
	if err != nil {
		// If Google OAuth verification fails, try Firebase token
		if os.Getenv("FIREBASE_AUTH_DOMAIN") != "" {
			// Parse the token without verification first to get the claims
			claims := jwt.MapClaims{}
			_, _, err := jwt.NewParser().ParseUnverified(req.IDToken, claims)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
				return
			}

			// Verify the token was issued by Firebase
			if aud, ok := claims["aud"].(string); !ok || aud != os.Getenv("FIREBASE_PROJECT_ID") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token audience"})
				return
			}

			// Verify the token was issued for our client
			if iss, ok := claims["iss"].(string); !ok || !strings.HasPrefix(iss, "https://securetoken.google.com/") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token issuer"})
				return
			}

			// Check if token is expired
			if exp, ok := claims["exp"].(float64); !ok || float64(time.Now().Unix()) > exp {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				return
			}

			// Get email from claims
			email, ok := claims["email"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
				return
			}

			tokenInfo = &oauth2api.Tokeninfo{
				Email: email,
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			return
		}
	}

	// Check if user exists
	var user models.User
	if err := database.DB.Where("email = ?", tokenInfo.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return 404 status to indicate user not found, so frontend can redirect to signup
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"email":   tokenInfo.Email,
				"action":  "signup",
				"message": "Please complete your registration",
			})
			return
		}
		// For any other database errors, return 500
		log.Printf("Database error in googleAuth: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Database error: %v", err)})
		return
	}

	// Generate JWT token
	tokenString := generateToken(user.ID, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}

func checkMobileNumber(c *gin.Context) {
	var req struct {
		MobileNumber string `json:"mobile_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mobile number is required"})
		return
	}

	var user models.User
	if err := database.DB.Where("mobile_number = ?", req.MobileNumber).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"exists": false,
				"message": "Mobile number not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists": true,
		"user_id": user.ID,
		"message": "Mobile number found",
	})
}

func oldLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Login request binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	log.Printf("Attempting login for email: %s", req.Email)
	log.Printf("Request body: %+v", req)

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User not found for email: %s", req.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		log.Printf("Database error during login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("User found in database: %+v", user)
	log.Printf("Stored password hash: %s", user.Password)
	log.Printf("Attempting to compare password hash with provided password")

	// Compare password with hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Password mismatch for email: %s, error: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	log.Printf("Login successful for email: %s", req.Email)
	// Generate token
	token := generateToken(user.ID, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
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
	})
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

	// Trim any whitespace from the user ID
	userID = strings.TrimSpace(userID)

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

// Test endpoint to get all users
func getAllUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Create a slice to hold user data with passwords
	type UserWithPassword struct {
		ID              string         `json:"id"`
		Name            string         `json:"name"`
		Username        string         `json:"username"`
		Email           string         `json:"email"`
		Password        string         `json:"password"` // Include password for testing
		Role            models.UserRole `json:"role"`
		AvatarURL       string         `json:"avatar_url"`
		Bio             string         `json:"bio"`
		MediaGallery    datatypes.JSON `json:"media_gallery"`
		WalletBalance   float64        `json:"wallet_balance"`
		FollowerCount   int            `json:"follower_count"`
		InstagramHandle string         `json:"instagram_handle"`
		Verified        bool           `json:"verified"`
		Age             int            `json:"age"`
		Gender          string         `json:"gender"`
		Location        string         `json:"location"`
		Height          float64        `json:"height"`
		Work            string         `json:"work"`
		Education       string         `json:"education"`
		MobileNumber    string         `json:"mobile_number"`
		CreatedAt       time.Time      `json:"created_at"`
		UpdatedAt       time.Time      `json:"updated_at"`
	}

	var usersWithPasswords []UserWithPassword
	for _, user := range users {
		usersWithPasswords = append(usersWithPasswords, UserWithPassword{
			ID:              user.ID,
			Name:            user.Name,
			Username:        user.Username,
			Email:           user.Email,
			Password:        user.Password, // Include the hashed password
			Role:            user.Role,
			AvatarURL:       user.AvatarURL,
			Bio:             user.Bio,
			MediaGallery:    user.MediaGallery,
			WalletBalance:   user.WalletBalance,
			FollowerCount:   user.FollowerCount,
			InstagramHandle: user.InstagramHandle,
			Verified:        user.Verified,
			Age:             user.Age,
			Gender:          string(user.Gender),
			Location:        user.Location,
			Height:          user.Height,
			Work:            user.Work,
			Education:       user.Education,
			MobileNumber:    user.MobileNumber,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": usersWithPasswords,
		"count": len(usersWithPasswords),
	})
}
