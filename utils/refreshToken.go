package utils

import (
	"context"
	"fmt"
	"time"
	"github.com/HersheyPlus/go-auth/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



// TokenStore handles refresh token operations
type TokenStore struct {
    db *gorm.DB
}

// NewTokenStore creates a new token store instance
func NewTokenStore(db *gorm.DB) *TokenStore {
    return &TokenStore{db: db}
}

// StoreToken stores a refresh token in the database
func (s *TokenStore) StoreToken(ctx context.Context, userID uuid.UUID, tokenDetails *TokenDetails) error {
    refreshToken := models.RefreshToken{
        UserID:    userID,
        TokenUUID: tokenDetails.RefreshUuid,
        ExpiresAt: tokenDetails.RtExpires,
    }

    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Optionally delete existing refresh tokens for the user
        if err := tx.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
            return fmt.Errorf("failed to delete existing tokens: %w", err)
        }

        // Store new refresh token
        if err := tx.Create(&refreshToken).Error; err != nil {
            return fmt.Errorf("failed to store refresh token: %w", err)
        }

        return nil
    })
}

// DeleteToken removes a refresh token from the database
func (s *TokenStore) DeleteToken(ctx context.Context, tokenUUID string) error {
    result := s.db.WithContext(ctx).Where("token_uuid = ?", tokenUUID).Delete(&models.RefreshToken{})
    if result.Error != nil {
        return fmt.Errorf("failed to delete token: %w", result.Error)
    }
    return nil
}

// ValidateToken checks if a refresh token exists and is valid
func (s *TokenStore) ValidateToken(ctx context.Context, tokenUUID string) (uuid.UUID, error) {
    var token models.RefreshToken
    err := s.db.WithContext(ctx).Where(
        "token_uuid = ? AND expires_at > ?", 
        tokenUUID, 
        time.Now(),
    ).First(&token).Error

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return uuid.Nil, fmt.Errorf("token not found or expired")
        }
        return uuid.Nil, fmt.Errorf("failed to validate token: %w", err)
    }

    return token.UserID, nil
}

// CleanupExpiredTokens removes expired refresh tokens
func (s *TokenStore) CleanupExpiredTokens(ctx context.Context) error {
    return s.db.WithContext(ctx).
        Where("expires_at < ?", time.Now()).
        Delete(&models.RefreshToken{}).Error
}