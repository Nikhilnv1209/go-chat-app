package repository

import (
	"context"
	"time"

	"chat-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageReceiptRepository struct {
	DB *gorm.DB
}

func NewMessageReceiptRepository(db *gorm.DB) MessageReceiptRepository {
	return &messageReceiptRepository{DB: db}
}

// Create creates a new message receipt
func (r *messageReceiptRepository) Create(ctx context.Context, receipt *models.MessageReceipt) error {
	receipt.ID = uuid.New()
	receipt.CreatedAt = time.Now()
	receipt.UpdatedAt = time.Now()
	return r.DB.WithContext(ctx).Create(receipt).Error
}

// CreateBatch creates multiple message receipts efficiently
func (r *messageReceiptRepository) CreateBatch(ctx context.Context, receipts []*models.MessageReceipt) error {
	// Pre-populate ID and timestamp if GORM hooks don't cover it (we are explicit here)
	for _, receipt := range receipts {
		if receipt.ID == uuid.Nil {
			receipt.ID = uuid.New()
		}
		if receipt.CreatedAt.IsZero() {
			receipt.CreatedAt = time.Now()
		}
		if receipt.UpdatedAt.IsZero() {
			receipt.UpdatedAt = time.Now()
		}
	}
	return r.DB.WithContext(ctx).Create(&receipts).Error
}

// UpdateStatus updates the status of a message receipt for a specific user
func (r *messageReceiptRepository) UpdateStatus(ctx context.Context, messageID, userID uuid.UUID, status string) error {
	return r.DB.WithContext(ctx).Model(&models.MessageReceipt{}).
		Where("message_id = ? AND user_id = ?", messageID, userID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// FindByMessageID returns all receipts for a specific message
func (r *messageReceiptRepository) FindByMessageID(ctx context.Context, messageID uuid.UUID) ([]models.MessageReceipt, error) {
	var receipts []models.MessageReceipt
	err := r.DB.WithContext(ctx).Where("message_id = ?", messageID).Find(&receipts).Error
	return receipts, err
}

// FindUnreadCount returns the count of unread messages for a user
// (Messages with receipts in SENT or DELIVERED status)
func (r *messageReceiptRepository) FindUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.MessageReceipt{}).
		Where("user_id = ? AND status != ?", userID, "READ").
		Count(&count).Error
	return count, err
}
