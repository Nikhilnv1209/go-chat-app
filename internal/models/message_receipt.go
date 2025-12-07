package models

import (
	"time"
)

type MessageReceipt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MessageID uint      `gorm:"not null;index" json:"message_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Status    string    `gorm:"size:20;default:'SENT'" json:"status"` // SENT, DELIVERED, READ
	UpdatedAt time.Time `json:"updated_at"`
}
