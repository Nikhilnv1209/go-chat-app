package repository

import (
	"chat-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(msg *models.Message) error {
	return r.db.Create(msg).Error
}

func (r *messageRepository) FindByConversation(userID, targetID uuid.UUID, msgType string, limit, beforeID int) ([]models.Message, error) {
	query := r.db.Order("created_at DESC").Limit(limit)

	// Filter by conversation type
	if msgType == "DM" {
		// For DMs: messages between userID and targetID
		query = query.Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, targetID, targetID, userID,
		)
	} else if msgType == "GROUP" {
		// For Groups: messages where group_id = targetID
		query = query.Where("group_id = ?", targetID)
	}

	// Pagination support (beforeID not currently used as we're using int, will be ignored)
	// Note: The beforeID parameter is int but IDs are UUIDs, so pagination via beforeID won't work
	// This is a known limitation for now

	var messages []models.Message
	err := query.Find(&messages).Error
	return messages, err
}
