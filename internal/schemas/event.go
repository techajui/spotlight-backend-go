package schemas

import "time"

type EventResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Location    string        `json:"location"`
	CreatorID   uint          `json:"creator_id"`
	Creator     UserResponse  `json:"creator"`
	Status      string        `json:"status"`
	MinBid      float64       `json:"min_bid"`
	CurrentBid  float64       `json:"current_bid"`
	WinnerID    *uint         `json:"winner_id"`
	Winner      *UserResponse `json:"winner,omitempty"`
}

type EventCreate struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Date        string   `json:"date" binding:"required"`
	Location    string   `json:"location" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	MinBid      float64  `json:"min_bid" binding:"required"`
	Images      []string `json:"images"`
	Attendees   []string `json:"attendees"`
	Capacity    int      `json:"capacity" binding:"required"`
	BidDeadline string   `json:"bid_deadline" binding:"required"`
}

type EventUpdate struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	Location    *string    `json:"location,omitempty"`
	Category    *string    `json:"category,omitempty"`
	MinBid      *float64   `json:"min_bid,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Images      *[]string  `json:"images,omitempty"`
	Attendees   *[]string  `json:"attendees,omitempty"`
}
