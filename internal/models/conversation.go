package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	BaseModel
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Type          string    `gorm:"size:10" json:"type"` // DM or GROUP
	TargetID      uuid.UUID `gorm:"type:uuid;not null" json:"target_id"`
	LastMessageAt time.Time `json:"last_message_at"`
	UnreadCount   int       `gorm:"default:0" json:"unread_count"`
}
