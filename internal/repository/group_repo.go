package repository

import (
	"context"
	"time"

	"chat-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) Create(ctx context.Context, group *models.Group) error {
	group.ID = uuid.New()
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *groupRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.WithContext(ctx).First(&group, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepository) GetMembers(ctx context.Context, groupID uuid.UUID) ([]models.GroupMember, error) {
	var members []models.GroupMember
	err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}

func (r *groupRepository) IsMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *groupRepository) AddMember(ctx context.Context, groupID, userID uuid.UUID, role string) error {
	member := &models.GroupMember{
		GroupID:  groupID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(member).Error
}
