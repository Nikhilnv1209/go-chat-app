package repository

import (
	"time"

	"chat-app/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	UpdateOnlineStatus(userID uuid.UUID, isOnline bool, lastSeen time.Time) error
}

type MessageRepository interface {
	Create(msg *models.Message) error
	FindByID(id uuid.UUID) (*models.Message, error)
	FindByConversation(userID, targetID uuid.UUID, msgType string, limit int, beforeID *uuid.UUID) ([]models.Message, error)
}

type MessageReceiptRepository interface {
	Create(receipt *models.MessageReceipt) error
	CreateBatch(receipts []*models.MessageReceipt) error
	UpdateStatus(messageID, userID uuid.UUID, status string) error
	FindByMessageID(messageID uuid.UUID) ([]models.MessageReceipt, error)
	FindUnreadCount(userID uuid.UUID) (int64, error)
}

type GroupRepository interface {
	Create(group *models.Group) error
	FindByID(id uuid.UUID) (*models.Group, error)
	GetMembers(groupID uuid.UUID) ([]models.GroupMember, error)
	IsMember(groupID, userID uuid.UUID) (bool, error)
	AddMember(groupID, userID uuid.UUID, role string) error
}

type ConversationRepository interface {
	Upsert(conv *models.Conversation) error
	FindByUser(userID uuid.UUID) ([]models.Conversation, error)
	IncrementUnread(userID uuid.UUID, convType string, targetID uuid.UUID) error
	ResetUnread(userID uuid.UUID, convType string, targetID uuid.UUID) error
}
