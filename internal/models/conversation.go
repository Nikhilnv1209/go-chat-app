package models

import (
	"time"
)

type Conversation struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Type          string    `gorm:"size:10" json:"type"` // DM or GROUP
	TargetID      uint      `gorm:"not null" json:"target_id"`
	LastMessageAt time.Time `json:"last_message_at"`
	UnreadCount   int       `gorm:"default:0" json:"unread_count"`
}
