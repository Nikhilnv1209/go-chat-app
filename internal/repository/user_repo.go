package repository

import (
	"context"
	"errors"
	"time"

	"chat-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific NotFound error?
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateOnlineStatus(ctx context.Context, userID uuid.UUID, isOnline bool, lastSeen time.Time) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_online": isOnline,
		"last_seen": lastSeen,
	}).Error
}

func (r *userRepository) Search(ctx context.Context, query string, excludeUserID uuid.UUID) ([]models.User, error) {
	var users []models.User
	db := r.db.WithContext(ctx).Where("id != ?", excludeUserID)

	if query != "" {
		searchPattern := "%" + query + "%"
		db = db.Where("username LIKE ? OR email LIKE ?", searchPattern, searchPattern)
	} else {
		// If query is empty, limit results (e.g., top 20 recent users or random)
		// For now, let's just limit to 20 to avoid dumping the whole DB
		db = db.Limit(20)
	}

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
