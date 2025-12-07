package models

import (
	"github.com/google/uuid"
)

type Message struct {
	BaseModel
	SenderID   uuid.UUID  `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID *uuid.UUID `gorm:"type:uuid" json:"receiver_id,omitempty"` // Nullable (for Groups)
	GroupID    *uuid.UUID `gorm:"type:uuid" json:"group_id,omitempty"`    // Nullable (for DMs)
	Content    string     `gorm:"type:text" json:"content"`
	MsgType    string     `gorm:"size:20;default:'TEXT'" json:"msg_type"`

	// Associations
	Sender User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}
