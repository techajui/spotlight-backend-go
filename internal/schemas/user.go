package schemas

import (
	"spotlight-backend-go/internal/models"
)

type UserResponse struct {
	ID                string                `json:"id"`
	Name              string                `json:"name"`
	Username          string                `json:"username"`
	Email             string                `json:"email"`
	AvatarURL         string                `json:"avatar_url"`
	Bio               string                `json:"bio"`
	Role              string                `json:"role"`
	WalletBalance     float64               `json:"wallet_balance"`
	MediaGallery      []string              `json:"media_gallery"`
	ProfilePhotos     []string              `json:"profile_photos"`
	CoverPhotoURL     string                `json:"cover_photo_url,omitempty"`
	FollowerCount     int                   `json:"follower_count"`
	InstagramHandle   string                `json:"instagram_handle,omitempty"`
	EventsAttended    []string              `json:"events_attended"`
	EventsHosted      []string              `json:"events_hosted"`
	EventsHostedCount int                   `json:"events_hosted_count"`
	Age               int                   `json:"age"`
	Gender            models.Gender         `json:"gender,omitempty"`
	Location          string                `json:"location,omitempty"`
	Height            float64               `json:"height,omitempty"`
	Work              string                `json:"work,omitempty"`
	Education         string                `json:"education,omitempty"`
	EducationLevel    models.EducationLevel `json:"education_level,omitempty"`
	Drinking          models.DrinkingStatus `json:"drinking,omitempty"`
	Interests         []string              `json:"interests,omitempty"`
	IsVerified        bool                  `json:"is_verified"`
	GovernmentIDURL   string                `json:"government_id_url,omitempty"`
	VerifiedAt        *string               `json:"verified_at,omitempty"`
	MobileNumber      string                `json:"mobile_number,omitempty"`
	CreatedAt         string                `json:"created_at"`
	UpdatedAt         string                `json:"updated_at"`
}

type UserUpdate struct {
	Name            *string                `json:"name,omitempty"`
	Username        *string                `json:"username,omitempty"`
	AvatarURL       *string                `json:"avatar_url,omitempty"`
	Bio             *string                `json:"bio,omitempty"`
	MediaGallery    *[]string              `json:"media_gallery,omitempty"`
	ProfilePhotos   *[]string              `json:"profile_photos,omitempty"`
	CoverPhotoURL   *string                `json:"cover_photo_url,omitempty"`
	InstagramHandle *string                `json:"instagram_handle,omitempty"`
	Age             *int                   `json:"age,omitempty"`
	Gender          *models.Gender         `json:"gender,omitempty"`
	Location        *string                `json:"location,omitempty"`
	Height          *float64               `json:"height,omitempty"`
	Work            *string                `json:"work,omitempty"`
	Education       *string                `json:"education,omitempty"`
	EducationLevel  *models.EducationLevel `json:"education_level,omitempty"`
	Drinking        *models.DrinkingStatus `json:"drinking,omitempty"`
	Interests       *[]string              `json:"interests,omitempty"`
	GovernmentIDURL *string                `json:"government_id_url,omitempty"`
	MobileNumber    *string                `json:"mobile_number,omitempty"`
}

type UserCreate struct {
	Name            string                `json:"name" binding:"required"`
	Username        string                `json:"username" binding:"required"`
	Email           string                `json:"email" binding:"required,email"`
	Password        string                `json:"password" binding:"required,min=6"`
	Role            string                `json:"role" binding:"required"`
	Age             int                   `json:"age" binding:"required,min=18"`
	Gender          models.Gender         `json:"gender,omitempty"`
	Location        string                `json:"location,omitempty"`
	Height          float64               `json:"height,omitempty"`
	Work            string                `json:"work,omitempty"`
	Education       string                `json:"education,omitempty"`
	EducationLevel  models.EducationLevel `json:"education_level,omitempty"`
	Drinking        models.DrinkingStatus `json:"drinking,omitempty"`
	Interests       []string              `json:"interests,omitempty"`
	MediaGallery    []string              `json:"media_gallery,omitempty"`
	EventsAttended  []string              `json:"events_attended,omitempty"`
	EventsHosted    []string              `json:"events_hosted,omitempty"`
	GovernmentIDURL string                `json:"government_id_url,omitempty"`
	MobileNumber    string                `json:"mobile_number,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}
