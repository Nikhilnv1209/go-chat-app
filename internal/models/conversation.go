package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	BaseModel
	UserID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_conversation" json:"user_id"`
	Type          string    `gorm:"size:10;uniqueIndex:idx_user_conversation" json:"type"` // DM or GROUP
	TargetID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_conversation" json:"target_id"`
	LastMessageAt time.Time `json:"last_message_at"`
	UnreadCount   int       `gorm:"default:0" json:"unread_count"`
}

// Add composite unique index
func (Conversation) TableName() string {
	return "conversations"
}

// Setting the unique index via GORM tag is also possible but sometimes messy with composites.
// Using the struct tag on one field or a separate function is better.
// Actually, GORM allows `gorm:"uniqueIndex:idx_user_conv"` on multiple fields to create a composite one.

// Let's modify the struct tags instead to be cleaner.
