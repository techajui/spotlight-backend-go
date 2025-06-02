package api

import (
	"net/http"
	"spotlight-backend-go/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Chat represents a chat between users
type Chat struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	User1ID   uint           `json:"user1_id"`
	User2ID   uint           `json:"user2_id"`
	User1     models.User    `json:"user1" gorm:"foreignKey:User1ID"`
	User2     models.User    `json:"user2" gorm:"foreignKey:User2ID"`
	Messages  []Message      `json:"messages"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Message represents a message in a chat
type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ChatID    uint           `json:"chat_id"`
	SenderID  uint           `json:"sender_id"`
	Content   string         `json:"content"`
	Read      bool           `json:"read"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// RegisterChatRoutes registers chat-related routes
func RegisterChatRoutes(router *gin.RouterGroup, db *gorm.DB) {
	chatGroup := router.Group("/chats")
	{
		chatGroup.GET("", getChats(db))
		chatGroup.GET("/:id", getChat(db))
		chatGroup.POST("", createChat(db))
		chatGroup.POST("/:id/messages", sendMessage(db))
		chatGroup.PUT("/:id/messages/:messageId/read", markMessageAsRead(db))
	}
}

// getChats returns all chats for the current user
func getChats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var chats []Chat
		if err := db.Preload("User1").Preload("User2").Preload("Messages").
			Where("user1_id = ? OR user2_id = ?", userID, userID).
			Find(&chats).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
			return
		}
		c.JSON(http.StatusOK, chats)
	}
}

// getChat returns a specific chat
func getChat(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		chatID := c.Param("id")
		var chat Chat
		if err := db.Preload("User1").Preload("User2").Preload("Messages").
			Where("id = ? AND (user1_id = ? OR user2_id = ?)", chatID, userID, userID).
			First(&chat).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
			return
		}
		c.JSON(http.StatusOK, chat)
	}
}

// createChat creates a new chat
func createChat(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var req struct {
			User2ID uint `json:"user2_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if chat already exists
		var existingChat Chat
		if err := db.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
			userID, req.User2ID, req.User2ID, userID).First(&existingChat).Error; err == nil {
			c.JSON(http.StatusOK, existingChat)
			return
		}

		chat := Chat{
			User1ID: userID,
			User2ID: req.User2ID,
		}
		if err := db.Create(&chat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
			return
		}
		c.JSON(http.StatusCreated, chat)
	}
}

// sendMessage sends a message in a chat
func sendMessage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		chatID := c.Param("id")
		var req struct {
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify chat exists and user is a participant
		var chat Chat
		if err := db.Where("id = ? AND (user1_id = ? OR user2_id = ?)", chatID, userID, userID).
			First(&chat).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
			return
		}

		message := Message{
			ChatID:   chat.ID,
			SenderID: userID,
			Content:  req.Content,
		}
		if err := db.Create(&message).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
			return
		}
		c.JSON(http.StatusCreated, message)
	}
}

// markMessageAsRead marks a message as read
func markMessageAsRead(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		chatID := c.Param("id")
		messageID := c.Param("messageId")

		// Verify chat exists and user is a participant
		var chat Chat
		if err := db.Where("id = ? AND (user1_id = ? OR user2_id = ?)", chatID, userID, userID).
			First(&chat).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
			return
		}

		// Update message
		if err := db.Model(&Message{}).
			Where("id = ? AND chat_id = ?", messageID, chatID).
			Update("read", true).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
	}
}

// Add a helper function to convert string ID to uint
func stringToUint(s string) uint {
	id, _ := strconv.ParseUint(s, 10, 64)
	return uint(id)
}
