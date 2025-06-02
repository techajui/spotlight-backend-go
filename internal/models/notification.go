package models

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeBidPlaced     NotificationType = "bid_placed"
	NotificationTypeBidAccepted   NotificationType = "bid_accepted"
	NotificationTypeBidRejected   NotificationType = "bid_rejected"
	NotificationTypeEventReminder NotificationType = "event_reminder"
	NotificationTypeNewFollower   NotificationType = "new_follower"
	NotificationTypePayment       NotificationType = "payment"
)

// Notification represents a notification for a user
type Notification struct {
	ID        string           `json:"id" bson:"_id,omitempty"`
	UserID    string           `json:"userId" bson:"user_id"`
	Type      NotificationType `json:"type" bson:"type"`
	Title     string           `json:"title" bson:"title"`
	Message   string           `json:"message" bson:"message"`
	Read      bool             `json:"read" bson:"read"`
	CreatedAt time.Time        `json:"createdAt" bson:"created_at"`
	// Optional fields for related entities
	RelatedEventID *string `json:"relatedEventId,omitempty" bson:"related_event_id,omitempty"`
	RelatedUserID  *string `json:"relatedUserId,omitempty" bson:"related_user_id,omitempty"`
	RelatedBidID   *string `json:"relatedBidId,omitempty" bson:"related_bid_id,omitempty"`
}
