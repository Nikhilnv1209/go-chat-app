package repository

import (
	"time"

	"chat-app/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	UpdateOnlineStatus(userID uint, isOnline bool, lastSeen time.Time) error
}

type MessageRepository interface {
	Create(msg *models.Message) error
	FindByConversation(userID, targetID uint, msgType string, limit, beforeID int) ([]models.Message, error)
}

type MessageReceiptRepository interface {
	Create(receipt *models.MessageReceipt) error
	UpdateStatus(messageID, userID uint, status string) error
	FindUnreadCount(userID uint) (int, error)
}

type GroupRepository interface {
	Create(group *models.Group) error
	FindByID(id uint) (*models.Group, error)
	GetMembers(groupID uint) ([]models.GroupMember, error)
	IsMember(groupID, userID uint) (bool, error)
	AddMember(groupID, userID uint, role string) error
}

type ConversationRepository interface {
	Upsert(conv *models.Conversation) error
	FindByUser(userID uint) ([]models.Conversation, error)
	IncrementUnread(userID uint, convType string, targetID uint) error
	ResetUnread(userID uint, convType string, targetID uint) error
}
