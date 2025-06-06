package views

import (
	"time"

	"github.com/trendscout/backend/internal/models"
)

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// NewUserResponse creates a new user response from a user model
func NewUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

// UserListResponse represents a list of users
type UserListResponse struct {
	Users []*UserResponse `json:"users"`
	Count int             `json:"count"`
}

// NewUserListResponse creates a new user list response
func NewUserListResponse(users []*models.User) *UserListResponse {
	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = NewUserResponse(user)
	}

	return &UserListResponse{
		Users: userResponses,
		Count: len(userResponses),
	}
} 