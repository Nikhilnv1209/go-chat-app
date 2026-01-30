package repository

import (
	"context"
	"chat-app/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, msg *models.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *messageRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	var msg models.Message
	err := r.db.WithContext(ctx).First(&msg, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepository) FindByConversation(ctx context.Context, userID, targetID uuid.UUID, msgType string, limit int, beforeID *uuid.UUID) ([]models.Message, error) {
	query := r.db.WithContext(ctx).Preload("Sender").Order("created_at DESC").Limit(limit)

	// Cursor-based pagination: if beforeID is provided, fetch messages older than that message
	if beforeID != nil {
		// Get the timestamp of the cursor message
		var cursorTime time.Time
		err := r.db.WithContext(ctx).Model(&models.Message{}).
			Select("created_at").
			Where("id = ?", beforeID).
			Scan(&cursorTime).Error
		if err != nil {
			return nil, err
		}
		query = query.Where("created_at < ?", cursorTime)
	}

	// Filter by conversation type
	switch msgType {
	case "DM":
		// For DMs: messages between userID and targetID
		query = query.Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, targetID, targetID, userID,
		)
	case "GROUP":
		// For Groups: messages where group_id = targetID
		query = query.Where("group_id = ?", targetID)
	}

	var messages []models.Message
	err := query.Find(&messages).Error
	return messages, err
}
