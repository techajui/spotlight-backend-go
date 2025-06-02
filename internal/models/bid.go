package models

import (
	"time"
)

// Bid represents a bid placed on an event
type Bid struct {
	ID        string    `json:"id" gorm:"primaryKey;type:char(36)"`
	EventID   string    `json:"event_id" gorm:"type:char(36)"`
	UserID    string    `json:"user_id" gorm:"type:char(36)"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	Event *Event `json:"event,omitempty" gorm:"foreignKey:EventID"`
	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
