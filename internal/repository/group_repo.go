package repository

import (
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

func (r *groupRepository) Create(group *models.Group) error {
	group.ID = uuid.New()
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	return r.db.Create(group).Error
}

func (r *groupRepository) FindByID(id uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.First(&group, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepository) GetMembers(groupID uuid.UUID) ([]models.GroupMember, error) {
	var members []models.GroupMember
	err := r.db.Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}

func (r *groupRepository) IsMember(groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *groupRepository) AddMember(groupID, userID uuid.UUID, role string) error {
	member := &models.GroupMember{
		GroupID:  groupID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
	return r.db.Create(member).Error
}
