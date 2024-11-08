package models

import (
	"gorm.io/gorm"
	"time"
	"github.com/google/uuid"
)

type Base struct {
	UserID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	LastLogin time.Time `gorm:"not null;default:current_timestamp"`
}

// Set UUID
func (base *Base) BeforeCreate() error {
	if base.UserID == uuid.Nil {
		base.UserID = uuid.New()
	}
	base.CreatedAt = time.Now()
	base.UpdatedAt = time.Now()
	return nil
}

// Set Time When Updated
func (base *Base) BeforeUpdate() error {
	base.UpdatedAt = time.Now()
	return nil
}