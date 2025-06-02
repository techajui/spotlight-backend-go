package models

import (
	"time"

	"gorm.io/datatypes"
)

// EventStatus represents the status of an event
type EventStatus string

const (
	EventStatusUpcoming  EventStatus = "upcoming"
	EventStatusOngoing   EventStatus = "ongoing"
	EventStatusPast      EventStatus = "past"
	EventStatusCancelled EventStatus = "cancelled"
)

// ApplicationStatus represents the status of an event application
type ApplicationStatus string

const (
	ApplicationStatusPending  ApplicationStatus = "pending"
	ApplicationStatusAccepted ApplicationStatus = "accepted"
	ApplicationStatusRejected ApplicationStatus = "rejected"
)

// Application represents a fan's application to attend an event
type Application struct {
	ID        string            `json:"id" bson:"_id,omitempty"`
	EventID   string            `json:"eventId" bson:"event_id"`
	FanID     string            `json:"fanId" bson:"fan_id"`
	BidAmount float64           `json:"bidAmount" bson:"bid_amount"`
	Message   string            `json:"message" bson:"message"`
	Status    ApplicationStatus `json:"status" bson:"status"`
	CreatedAt time.Time         `json:"createdAt" bson:"created_at"`
}

// Event represents an event hosted by an influencer
type Event struct {
	ID          string         `json:"id" gorm:"primaryKey;type:char(36)"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Date        time.Time      `json:"date"`
	Location    string         `json:"location"`
	HostID      string         `json:"host_id" gorm:"type:char(36)"`
	Category    string         `json:"category"`
	Images      datatypes.JSON `json:"images"`
	MinBid      float64        `json:"min_bid"`
	Capacity    int            `json:"capacity"`
	BidDeadline time.Time      `json:"bid_deadline"`
	Status      EventStatus    `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Associations
	Host      *User  `json:"host,omitempty" gorm:"foreignKey:HostID"`
	Attendees []User `json:"attendees,omitempty" gorm:"many2many:event_attendees;"`
}
