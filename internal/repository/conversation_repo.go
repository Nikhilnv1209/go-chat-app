package repository

import (
	"chat-app/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) Upsert(conv *models.Conversation) error {
	// Upsert: On conflict (user_id, type, target_id), update last_message_at
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "type"}, {Name: "target_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_message_at"}),
	}).Create(conv).Error
}

func (r *conversationRepository) FindByUser(userID uuid.UUID) ([]models.Conversation, error) {
	var convs []models.Conversation
	err := r.db.Where("user_id = ?", userID).Order("last_message_at DESC").Find(&convs).Error
	return convs, err
}

func (r *conversationRepository) IncrementUnread(userID uuid.UUID, convType string, targetID uuid.UUID) error {
	// We need to upsert first to ensure the row exists, but GORM's atomic increment is tricky with Upsert in one go.
	// For MVP, we can try to update, if 0 rows affected, then create.

	result := r.db.Model(&models.Conversation{}).
		Where("user_id = ? AND type = ? AND target_id = ?", userID, convType, targetID).
		Update("unread_count", gorm.Expr("unread_count + ?", 1))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// Create new if not exists
		newConv := models.Conversation{
			UserID:        userID,
			Type:          convType,
			TargetID:      targetID,
			LastMessageAt: time.Now(),
			UnreadCount:   1,
		}
		return r.db.Create(&newConv).Error
	}

	// Also update timestamp
	return r.db.Model(&models.Conversation{}).
		Where("user_id = ? AND type = ? AND target_id = ?", userID, convType, targetID).
		Update("last_message_at", time.Now()).Error
}

func (r *conversationRepository) ResetUnread(userID uuid.UUID, convType string, targetID uuid.UUID) error {
	return r.db.Model(&models.Conversation{}).
		Where("user_id = ? AND type = ? AND target_id = ?", userID, convType, targetID).
		Update("unread_count", 0).Error
}
