package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;unique;not null" json:"username"`
	Email     string    `gorm:"size:100;unique;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"` // Never result password
	IsOnline  bool      `gorm:"default:false" json:"is_online"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
