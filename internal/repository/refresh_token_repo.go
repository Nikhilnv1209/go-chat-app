package repository

import (
	"chat-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) GetByHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Where("token_hash = ?", hash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepository) Revoke(id uuid.UUID) error {
	return r.db.Model(&models.RefreshToken{}).Where("id = ?", id).Update("revoked", true).Error
}

func (r *refreshTokenRepository) RevokeByUser(userID uuid.UUID) error {
	return r.db.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}
