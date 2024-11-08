package dto

import (
	"github.com/google/uuid"
	"time"
)

type UserRegisterResponse struct {
    Username string    `json:"username"`
    FirstName *string   `json:"first_name,omitempty"`
    LastName  *string   `json:"last_name,omitempty"`
    Phone     string    `json:"phone"`
    Email    string    `json:"email"`
    LastLogin time.Time `json:"last_login,omitempty"`
}

type UserLoginResponse struct {
    UserID           uuid.UUID `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token,omitempty"`
    ExpiresIn    int64     `json:"expires_in"`
    LastLogin time.Time `json:"last_login,omitempty"`
}



type UserProfileResponse struct {
    UserID        uuid.UUID `json:"id"`
    Username  string    `json:"username"`
    FirstName *string   `json:"first_name,omitempty"`
    LastName  *string   `json:"last_name,omitempty"`
    Phone     string    `json:"phone"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}