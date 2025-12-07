package models

import (
	"time"
	"gorm.io/gorm"
)

type Message struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	SenderID        uint           `gorm:"not null" json:"sender_id"`
	ReceiverID      *uint          `json:"receiver_id,omitempty"` // Nullable (for Groups)
	GroupID         *uint          `json:"group_id,omitempty"`    // Nullable (for DMs)
	Content         string         `gorm:"type:text" json:"content"`
	MsgType         string         `gorm:"size:20;default:'TEXT'" json:"msg_type"`
	CreatedAt       time.Time      `json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Associations
	Sender   User   `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}
