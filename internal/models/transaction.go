package models

import (
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeBidPlaced       TransactionType = "bid_placed"
	TransactionTypeBidRejected     TransactionType = "bid_rejected"
	TransactionTypePaymentReceived TransactionType = "payment_received"
	TransactionTypeTopUp           TransactionType = "top_up"
	TransactionTypeWithdrawal      TransactionType = "withdrawal"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID             string          `json:"id" bson:"_id,omitempty"`
	UserID         string          `json:"userId" bson:"user_id"`
	Type           TransactionType `json:"type" bson:"type"`
	Amount         float64         `json:"amount" bson:"amount"`
	Description    string          `json:"description" bson:"description"`
	CreatedAt      time.Time       `json:"createdAt" bson:"created_at"`
	RelatedEventID string          `json:"relatedEventId,omitempty" bson:"related_event_id,omitempty"`
	RelatedUserID  string          `json:"relatedUserId,omitempty" bson:"related_user_id,omitempty"`
}
