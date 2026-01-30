package service

import (
	"context"
	"chat-app/internal/models"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (string, string, *models.User, error)
	Login(ctx context.Context, email, password string) (string, string, *models.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateToken(tokenString string) (uuid.UUID, error)
	SearchUsers(ctx context.Context, query string, excludeUserID uuid.UUID) ([]models.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type MessageService interface {
	SendDirectMessage(ctx context.Context, senderID, receiverID uuid.UUID, content string) (*models.Message, error)
	SendGroupMessage(ctx context.Context, senderID, groupID uuid.UUID, content string) (*models.Message, error)
	GetHistory(ctx context.Context, userID, targetID uuid.UUID, convType string, limit int, beforeID *uuid.UUID) ([]models.Message, error)
	MarkAsRead(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) error
	MarkAsDelivered(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) error
	GetMessageReceipts(ctx context.Context, userID, messageID uuid.UUID) ([]models.MessageReceipt, error)
	GetUserInfo(ctx context.Context, userID uuid.UUID) (*models.User, error)
	BroadcastTypingIndicator(ctx context.Context, userID uuid.UUID, username, convType string, targetID uuid.UUID, isTyping bool) error
}

type GroupService interface {
	Create(ctx context.Context, creatorID uuid.UUID, name string, memberIDs []uuid.UUID) (*models.Group, error)
	AddMember(ctx context.Context, adminID, groupID, newMemberID uuid.UUID) error
	RemoveMember(ctx context.Context, adminID, groupID, memberID uuid.UUID) error
}
