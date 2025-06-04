package api

import (
	"encoding/json"
	"log"
	"net/http"
	"spotlight-backend-go/internal/database"
	"spotlight-backend-go/internal/models"
	"spotlight-backend-go/internal/schemas"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("/me", getMe)
		users.PUT("/me", updateMe)
		users.GET("/:id", getUserProfile)
		users.GET("/influencers", getAllInfluencers)
		users.GET("/me/events", getUserEvents)
	}
}

func getMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	var mediaGallery []string
	_ = json.Unmarshal(currentUser.MediaGallery, &mediaGallery)

	// Get hosted events for influencers
	var hostedEvents []models.Event
	if currentUser.Role == models.RoleInfluencer {
		database.DB.Model(currentUser).Association("HostedEvents").Find(&hostedEvents)
	}

	// Get attended events
	var attendedEvents []models.Event
	database.DB.Model(currentUser).Association("AttendedEvents").Find(&attendedEvents)

	log.Printf("Fetched avatar for user %s: %s", currentUser.ID, currentUser.AvatarURL)

	response := gin.H{
		"user":                currentUser,
		"media_gallery":       mediaGallery,
		"events_attended":     attendedEvents,
		"events_hosted":       hostedEvents,
		"events_hosted_count": currentUser.EventsHostedCount,
	}
	c.JSON(http.StatusOK, response)
}

func updateMe(c *gin.Context) {
	var updateData schemas.UserUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	if updateData.Name != nil {
		currentUser.Name = *updateData.Name
	}
	if updateData.AvatarURL != nil {
		log.Printf("Updating avatar for user %s: %s", currentUser.ID, *updateData.AvatarURL)
		currentUser.AvatarURL = *updateData.AvatarURL
	}
	if updateData.Bio != nil {
		currentUser.Bio = *updateData.Bio
	}
	if updateData.MediaGallery != nil {
		mediaGalleryJSON, _ := json.Marshal(*updateData.MediaGallery)
		currentUser.MediaGallery = datatypes.JSON(mediaGalleryJSON)
	}
	if updateData.ProfilePhotos != nil {
		profilePhotosJSON, _ := json.Marshal(*updateData.ProfilePhotos)
		currentUser.ProfilePhotos = datatypes.JSON(profilePhotosJSON)
		// Update avatar_url to be the first profile photo if available
		if len(*updateData.ProfilePhotos) > 0 {
			currentUser.AvatarURL = (*updateData.ProfilePhotos)[0]
		}
	}
	if updateData.Age != nil {
		currentUser.Age = *updateData.Age
	}
	if updateData.Gender != nil {
		currentUser.Gender = *updateData.Gender
	}
	if updateData.Location != nil {
		currentUser.Location = *updateData.Location
	}
	if updateData.Height != nil {
		currentUser.Height = *updateData.Height
	}
	if updateData.Work != nil {
		currentUser.Work = *updateData.Work
	}
	if updateData.Education != nil {
		currentUser.Education = *updateData.Education
	}
	if updateData.EducationLevel != nil {
		currentUser.EducationLevel = *updateData.EducationLevel
	}
	if updateData.Drinking != nil {
		currentUser.Drinking = *updateData.Drinking
	}
	if updateData.Interests != nil {
		interestsJSON, _ := json.Marshal(*updateData.Interests)
		currentUser.Interests = datatypes.JSON(interestsJSON)
	}
	if updateData.GovernmentIDURL != nil {
		currentUser.GovernmentIDURL = *updateData.GovernmentIDURL
	}

	// Update influencer-specific fields if user is an influencer
	if currentUser.Role == models.RoleInfluencer {
		if updateData.CoverPhotoURL != nil {
			currentUser.CoverPhotoURL = *updateData.CoverPhotoURL
		}
		if updateData.InstagramHandle != nil {
			currentUser.InstagramHandle = *updateData.InstagramHandle
		}
	}

	if err := database.DB.Save(currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	var mediaGallery []string
	_ = json.Unmarshal(currentUser.MediaGallery, &mediaGallery)

	// Get hosted events for influencers
	var hostedEvents []models.Event
	if currentUser.Role == models.RoleInfluencer {
		database.DB.Model(currentUser).Association("HostedEvents").Find(&hostedEvents)
	}

	// Get attended events
	var attendedEvents []models.Event
	database.DB.Model(currentUser).Association("AttendedEvents").Find(&attendedEvents)

	response := gin.H{
		"user":                currentUser,
		"media_gallery":       mediaGallery,
		"events_attended":     attendedEvents,
		"events_hosted":       hostedEvents,
		"events_hosted_count": currentUser.EventsHostedCount,
	}
	c.JSON(http.StatusOK, response)
}

func getUserProfile(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new user if not found
			user = models.User{
				ID:        id,
				Name:      "New User",
				Username:  "user_" + id,
				Role:      models.RoleFan,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := database.DB.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			return
		}
	}

	var mediaGallery []string
	_ = json.Unmarshal(user.MediaGallery, &mediaGallery)

	// Get hosted events for influencers
	var hostedEvents []models.Event
	if user.Role == models.RoleInfluencer {
		database.DB.Model(&user).Association("HostedEvents").Find(&hostedEvents)
	}

	// Get attended events
	var attendedEvents []models.Event
	database.DB.Model(&user).Association("AttendedEvents").Find(&attendedEvents)

	response := gin.H{
		"user":                user,
		"media_gallery":       mediaGallery,
		"events_attended":     attendedEvents,
		"events_hosted":       hostedEvents,
		"events_hosted_count": user.EventsHostedCount,
	}
	c.JSON(http.StatusOK, response)
}

func getAllInfluencers(c *gin.Context) {
	var influencers []models.User
	if err := database.DB.Where("role = ?", models.RoleInfluencer).Find(&influencers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch influencers"})
		return
	}
	c.JSON(http.StatusOK, influencers)
}

// getUserEvents returns all events for the current user with bid status
func getUserEvents(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	// Get all events where user has placed a bid
	var events []models.Event
	if err := database.DB.
		Joins("JOIN bids ON bids.event_id = events.id").
		Where("bids.user_id = ?", currentUser.ID).
		Preload("Host").
		Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	// Get bid information for each event
	type EventWithBid struct {
		models.Event
		BidAmount float64 `json:"bid_amount"`
		BidStatus string  `json:"bid_status"`
	}

	var eventsWithBids []EventWithBid
	for _, event := range events {
		var bid models.Bid
		if err := database.DB.Where("event_id = ? AND user_id = ?", event.ID, currentUser.ID).First(&bid).Error; err != nil {
			continue
		}

		// Determine bid status
		bidStatus := "pending"
		if event.BidDeadline.Before(time.Now()) {
			bidStatus = "expired"
		}

		eventsWithBids = append(eventsWithBids, EventWithBid{
			Event:     event,
			BidAmount: bid.Amount,
			BidStatus: bidStatus,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"events": eventsWithBids,
	})
}
