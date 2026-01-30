package repository

import (
	"context"
	"time"

	"chat-app/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateOnlineStatus(ctx context.Context, userID uuid.UUID, isOnline bool, lastSeen time.Time) error
	Search(ctx context.Context, query string, excludeUserID uuid.UUID) ([]models.User, error)
}

type MessageRepository interface {
	Create(ctx context.Context, msg *models.Message) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	FindByConversation(ctx context.Context, userID, targetID uuid.UUID, msgType string, limit int, beforeID *uuid.UUID) ([]models.Message, error)
}

type MessageReceiptRepository interface {
	Create(ctx context.Context, receipt *models.MessageReceipt) error
	CreateBatch(ctx context.Context, receipts []*models.MessageReceipt) error
	UpdateStatus(ctx context.Context, messageID, userID uuid.UUID, status string) error
	FindByMessageID(ctx context.Context, messageID uuid.UUID) ([]models.MessageReceipt, error)
	FindUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)
}

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Group, error)
	GetMembers(ctx context.Context, groupID uuid.UUID) ([]models.GroupMember, error)
	IsMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error)
	AddMember(ctx context.Context, groupID, userID uuid.UUID, role string) error
}

type ConversationRepository interface {
	Upsert(ctx context.Context, conv *models.Conversation) error
	FindByUser(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error)
	IncrementUnread(ctx context.Context, userID uuid.UUID, convType string, targetID uuid.UUID, lastMessage string) error
	ResetUnread(ctx context.Context, userID uuid.UUID, convType string, targetID uuid.UUID) error
	FindContactsOfUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) // For presence broadcasting
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*models.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeByUser(ctx context.Context, userID uuid.UUID) error
}
