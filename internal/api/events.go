package api

import (
	"net/http"
	"spotlight-backend-go/internal/models"
	"spotlight-backend-go/internal/schemas"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterEventRoutes registers event-related routes
func RegisterEventRoutes(router *gin.RouterGroup, db *gorm.DB) {
	eventGroup := router.Group("/events")
	{
		eventGroup.GET("", getEvents(db))
		eventGroup.GET("/:id", getEvent(db))
		eventGroup.POST("", createEvent(db))
		eventGroup.PUT("/:id", updateEvent(db))
		eventGroup.DELETE("/:id", deleteEvent(db))
		eventGroup.POST("/:id/attend", attendEvent(db))
		eventGroup.POST("/:id/unattend", unattendEvent(db))
		eventGroup.POST("/:id/bid", placeBid(db))
	}
}

// getEvents returns all events
func getEvents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var events []models.Event
		if err := db.Preload("Host").Preload("Attendees").Find(&events).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
			return
		}
		c.JSON(http.StatusOK, events)
	}
}

// getEvent returns a specific event
func getEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var event models.Event
		if err := db.Preload("Host").Preload("Attendees").Where("id = ?", id).First(&event).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}
		c.JSON(http.StatusOK, event)
	}
}

// createEvent creates a new event
func createEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req schemas.EventCreate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse date strings into time.Time objects
		date, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use ISO 8601 format (e.g., 2024-03-20T15:00:00Z)"})
			return
		}

		bidDeadline, err := time.Parse(time.RFC3339, req.BidDeadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid deadline format. Use ISO 8601 format (e.g., 2024-03-20T15:00:00Z)"})
			return
		}

		userID := c.GetString("user_id")
		if userID == "" {
			userID = "1" // Default to a known influencer for local/dev
		}

		// Generate a new UUID for the event
		eventID := uuid.New().String()

		event := models.Event{
			ID:          eventID,
			Title:       req.Title,
			Description: req.Description,
			Date:        date,
			Location:    req.Location,
			HostID:      userID,
			Category:    req.Category,
			Images:      ToJSON(req.Images),
			MinBid:      req.MinBid,
			Capacity:    req.Capacity,
			BidDeadline: bidDeadline,
			Status:      models.EventStatusUpcoming,
		}

		if err := db.Create(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
			return
		}
		c.JSON(http.StatusCreated, event)
	}
}

// updateEvent updates an existing event
func updateEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req schemas.EventUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := c.GetString("user_id")
		var event models.Event
		if err := db.First(&event, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		// Check if user is the host
		if event.HostID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the host can update this event"})
			return
		}

		// Update fields
		if req.Title != nil {
			event.Title = *req.Title
		}
		if req.Description != nil {
			event.Description = *req.Description
		}
		if req.Date != nil {
			event.Date = *req.Date
		}
		if req.Location != nil {
			event.Location = *req.Location
		}
		if req.MinBid != nil {
			event.MinBid = *req.MinBid
		}
		if req.Images != nil {
			event.Images = ToJSON(*req.Images)
		}
		if req.Status != nil {
			event.Status = models.EventStatus(*req.Status)
		}

		if err := db.Save(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
			return
		}
		c.JSON(http.StatusOK, event)
	}
}

// deleteEvent deletes an event
func deleteEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.GetString("user_id")

		var event models.Event
		if err := db.First(&event, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		// Check if user is the host
		if event.HostID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the host can delete this event"})
			return
		}

		if err := db.Delete(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
	}
}

// attendEvent allows a user to attend an event
func attendEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.GetString("user_id")

		var event models.Event
		if err := db.First(&event, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		// Check if event is full
		var attendeeCount int64
		db.Model(&models.EventAttendee{}).Where("event_id = ?", event.ID).Count(&attendeeCount)
		if attendeeCount >= int64(event.Capacity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Event is full"})
			return
		}

		// Add attendee
		attendee := models.EventAttendee{
			EventID: event.ID,
			UserID:  userID,
		}
		if err := db.Create(&attendee).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to attend event"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successfully attended event"})
	}
}

// unattendEvent allows a user to unattend an event
func unattendEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.GetString("user_id")

		var event models.Event
		if err := db.First(&event, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		// Remove attendee
		if err := db.Where("event_id = ? AND user_id = ?", event.ID, userID).
			Delete(&models.EventAttendee{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unattend event"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successfully unattended event"})
	}
}

// placeBid allows a user to place a bid on an event
func placeBid(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.GetString("user_id")

		var req struct {
			Amount float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var event models.Event
		if err := db.First(&event, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		// Check if bid deadline has passed
		if event.BidDeadline.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bid deadline has passed"})
			return
		}

		// Check if bid amount is valid
		if req.Amount < event.MinBid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bid amount must be greater than minimum bid"})
			return
		}

		// Create bid
		bid := models.Bid{
			EventID: event.ID,
			UserID:  userID,
			Amount:  req.Amount,
		}
		if err := db.Create(&bid).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place bid"})
			return
		}
		c.JSON(http.StatusCreated, bid)
	}
}
