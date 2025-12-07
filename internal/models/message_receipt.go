package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageReceipt struct {
	BaseModel
	MessageID uuid.UUID `gorm:"type:uuid;not null;index" json:"message_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Status    string    `gorm:"size:20;default:'SENT'" json:"status"` // SENT, DELIVERED, READ
	UpdatedAt time.Time `json:"updated_at"`
}
