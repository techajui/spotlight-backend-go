package schemas

import "time"

type ChatRoomResponse struct {
	ID        uint          `json:"id"`
	EventID   uint          `json:"event_id"`
	Event     EventResponse `json:"event"`
	Status    string        `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type MessageResponse struct {
	ID         uint         `json:"id"`
	ChatRoomID uint         `json:"chat_room_id"`
	SenderID   uint         `json:"sender_id"`
	Sender     UserResponse `json:"sender"`
	Content    string       `json:"content"`
	ReadAt     *time.Time   `json:"read_at"`
	CreatedAt  time.Time    `json:"created_at"`
}

type MessageCreate struct {
	Content string `json:"content" binding:"required"`
}

type ChatRoomCreate struct {
	EventID uint `json:"event_id" binding:"required"`
}
