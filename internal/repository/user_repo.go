package repository

import (
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

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific NotFound error?
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateOnlineStatus(userID uuid.UUID, isOnline bool, lastSeen time.Time) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_online": isOnline,
		"last_seen": lastSeen,
	}).Error
}

func (r *userRepository) Search(query string, excludeUserID uuid.UUID) ([]models.User, error) {
	var users []models.User
	db := r.db.Where("id != ?", excludeUserID)

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
