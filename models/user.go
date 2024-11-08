package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    Base
    Username  string  `gorm:"type:varchar(100);not null;index" json:"username"`
    FirstName *string `gorm:"type:varchar(100)" json:"first_name,omitempty"`
    LastName  *string `gorm:"type:varchar(100)" json:"last_name,omitempty"`
    Phone     string  `gorm:"type:varchar(20);not null" json:"phone"`
    Email     string  `gorm:"type:varchar(100);unique;not null;index" json:"email"`
    Password  string  `gorm:"type:varchar(255);not null" json:"-"`
    RefreshTokens []RefreshToken `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenUUID string    `gorm:"type:varchar(100);not null;unique;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"not null;default:current_timestamp"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
