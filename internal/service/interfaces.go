package service

import (
	"chat-app/internal/models"
)

type AuthService interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
	ValidateToken(tokenString string) (uint, error)
}

type MessageService interface {
	SendDirectMessage(senderID, receiverID uint, content string) (*models.Message, error)
	SendGroupMessage(senderID, groupID uint, content string) (*models.Message, error)
	GetHistory(userID, targetID uint, convType string, limit, beforeID int) ([]models.Message, error)
	MarkAsRead(userID uint, messageIDs []uint) error
}

type GroupService interface {
	Create(creatorID uint, name string, memberIDs []uint) (*models.Group, error)
	AddMember(adminID, groupID, newMemberID uint) error
	RemoveMember(adminID, groupID, memberID uint) error
}
