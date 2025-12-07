package models

import (
	"time"
)

type User struct {
	BaseModel
	Username string    `gorm:"size:50;unique;not null" json:"username"`
	Email    string    `gorm:"size:100;unique;not null" json:"email"`
	Password string    `gorm:"size:255;not null" json:"-"` // Never result password
	IsOnline bool      `gorm:"default:false" json:"is_online"`
	LastSeen time.Time `json:"last_seen"`
}
