package models

import (
	"time"

	"github.com/google/uuid"
)

type GroupMember struct {
	GroupID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"group_id"`
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Role     string    `gorm:"size:20;default:'MEMBER'" json:"role"` // ADMIN, MEMBER
	JoinedAt time.Time `json:"joined_at"`
}
