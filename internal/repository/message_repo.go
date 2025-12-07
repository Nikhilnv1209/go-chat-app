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
	// Simple implementation for now - improve query for complex conversation logic later
	// This assumes we want messages where (sender=userID AND receiver=targetID) OR (sender=targetID AND receiver=userID)
	// msgType would be 'direct' usually.

	query := r.db.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, targetID, targetID, userID,
	).Order("created_at DESC").Limit(limit)

	if beforeID > 0 {
		query = query.Where("id < ?", beforeID)
	}

	var messages []models.Message
	err := query.Find(&messages).Error
	return messages, err
}
