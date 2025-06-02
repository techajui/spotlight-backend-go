package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatRoom struct {
	gorm.Model
	EventID uint   `json:"event_id"`
	Event   Event  `json:"event" gorm:"foreignKey:EventID"`
	Status  string `json:"status" gorm:"default:'active'"`
}

type Message struct {
	gorm.Model
	ChatRoomID uint       `json:"chat_room_id"`
	ChatRoom   ChatRoom   `json:"chat_room" gorm:"foreignKey:ChatRoomID"`
	SenderID   string     `json:"sender_id"`
	Sender     User       `json:"sender" gorm:"foreignKey:SenderID;references:ID"`
	Content    string     `json:"content"`
	ReadAt     *time.Time `json:"read_at"`
}
