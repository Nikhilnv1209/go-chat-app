package models

import (
	"time"
)

type GroupMember struct {
	GroupID  uint      `gorm:"primaryKey" json:"group_id"`
	UserID   uint      `gorm:"primaryKey" json:"user_id"`
	Role     string    `gorm:"size:20;default:'MEMBER'" json:"role"` // ADMIN, MEMBER
	JoinedAt time.Time `json:"joined_at"`
}
