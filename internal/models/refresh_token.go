package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a long-lived session token
type RefreshToken struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"` // Store hash, not raw token
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	IPAddress string    `gorm:"size:45" json:"ip_address,omitempty"` // IPv6 can be 45 chars
	UserAgent string    `gorm:"size:255" json:"user_agent,omitempty"`
}
