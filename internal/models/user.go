package models

import (
	"time"

	"gorm.io/datatypes"
)

// UserRole represents the type of user
type UserRole string

const (
	RoleFan        UserRole = "fan"
	RoleInfluencer UserRole = "influencer"
)

// Gender represents the user's gender
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// EducationLevel represents the user's education level
type EducationLevel string

const (
	EducationHighSchool EducationLevel = "high_school"
	EducationBachelors  EducationLevel = "bachelors"
	EducationMasters    EducationLevel = "masters"
	EducationPhD        EducationLevel = "phd"
	EducationOther      EducationLevel = "other"
)

// DrinkingStatus represents the user's drinking preference
type DrinkingStatus string

const (
	DrinkingYes    DrinkingStatus = "yes"
	DrinkingNo     DrinkingStatus = "no"
	DrinkingSocial DrinkingStatus = "social"
	DrinkingRarely DrinkingStatus = "rarely"
)

type User struct {
	ID            string         `json:"id" gorm:"primaryKey;type:char(36)"`
	Name          string         `json:"name"`
	Username      string         `json:"username" gorm:"unique"`
	Email         string         `json:"email" gorm:"unique"`
	Password      string         `json:"-"`
	AvatarURL     string         `json:"avatar_url"`
	Bio           string         `json:"bio"`
	Role          UserRole       `json:"role"`
	WalletBalance float64        `json:"wallet_balance"`
	MediaGallery  datatypes.JSON `json:"media_gallery"`
	ProfilePhotos datatypes.JSON `json:"profile_photos"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`

	// Basic Profile Information
	Age            int            `json:"age" gorm:"not null"`
	Gender         Gender         `json:"gender"`
	Location       string         `json:"location"`
	Height         float64        `json:"height"` // in cm
	Work           string         `json:"work"`
	Education      string         `json:"education"`
	EducationLevel EducationLevel `json:"education_level"`
	Drinking       DrinkingStatus `json:"drinking"`
	MobileNumber   string         `json:"mobile_number"`
	Interests      datatypes.JSON `json:"interests"` // Array of interest categories

	// Document Verification
	GovernmentIDURL string     `json:"government_id_url"`
	IsVerified      bool       `json:"is_verified" gorm:"default:false"`
	VerifiedAt      *time.Time `json:"verified_at"`

	// Influencer specific fields
	CoverPhotoURL   string `json:"cover_photo_url,omitempty"`
	FollowerCount   int    `json:"follower_count,omitempty"`
	InstagramHandle string `json:"instagram_handle,omitempty"`
	Verified        bool   `json:"verified,omitempty"`

	// Events hosted count
	EventsHostedCount int `json:"events_hosted_count" gorm:"default:0"`

	// Associations
	HostedEvents   []Event `json:"hosted_events,omitempty" gorm:"foreignKey:HostID"`
	AttendedEvents []Event `json:"attended_events,omitempty" gorm:"many2many:event_attendees;"`
}

// EventAttendee represents the many-to-many relationship between users and events
type EventAttendee struct {
	UserID    string `gorm:"primaryKey;type:char(36)"`
	EventID   string `gorm:"primaryKey;type:char(36)"`
	CreatedAt time.Time
}

// TableName specifies the table name for EventAttendee
func (EventAttendee) TableName() string {
	return "event_attendees"
}
