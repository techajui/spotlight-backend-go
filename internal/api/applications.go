package api

import (
	"net/http"
	"spotlight-backend-go/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getApplicationsByEventID fetches applications for a specific event
func getApplicationsByEventID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := c.Param("eventId")
		var applications []models.Application
		if err := db.Where("event_id = ?", eventID).Find(&applications).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
			return
		}
		c.JSON(http.StatusOK, applications)
	}
}

func RegisterApplicationRoutes(router *gin.RouterGroup, db *gorm.DB) {
	applicationGroup := router.Group("/applications")
	{
		applicationGroup.GET("/event/:eventId", getApplicationsByEventID(db))
	}
}
