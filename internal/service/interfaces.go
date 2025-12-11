package service

import (
	"chat-app/internal/models"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}

type MessageService interface {
	SendDirectMessage(senderID, receiverID uuid.UUID, content string) (*models.Message, error)
	SendGroupMessage(senderID, groupID uuid.UUID, content string) (*models.Message, error)
	GetHistory(userID, targetID uuid.UUID, convType string, limit, beforeID int) ([]models.Message, error)
	MarkAsRead(userID uuid.UUID, messageIDs []uuid.UUID) error
	MarkAsDelivered(userID uuid.UUID, messageIDs []uuid.UUID) error
	GetMessageReceipts(userID, messageID uuid.UUID) ([]models.MessageReceipt, error)
	GetUserInfo(userID uuid.UUID) (*models.User, error)
	BroadcastTypingIndicator(userID uuid.UUID, username, convType string, targetID uuid.UUID, isTyping bool) error
}

type GroupService interface {
	Create(creatorID uuid.UUID, name string, memberIDs []uuid.UUID) (*models.Group, error)
	AddMember(adminID, groupID, newMemberID uuid.UUID) error
	RemoveMember(adminID, groupID, memberID uuid.UUID) error
}
