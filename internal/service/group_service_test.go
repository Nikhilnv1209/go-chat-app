package service_test

import (
	"context"
	"chat-app/internal/models"
	"chat-app/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGroupService_Create_Success(t *testing.T) {
	ctx := context.Background()
	mockGroupRepo := new(MockGroupRepo)
	svc := service.NewGroupService(mockGroupRepo)

	creatorID := uuid.New()
	memberIDs := []uuid.UUID{uuid.New(), uuid.New()}
	groupName := "Test Group"

	// Mock group creation
	mockGroupRepo.On("Create", ctx, mock.MatchedBy(func(g *models.Group) bool {
		return g.Name == groupName
	})).Return(nil)

	// Mock adding creator as admin
	mockGroupRepo.On("AddMember", ctx, mock.Anything, creatorID, "ADMIN").Return(nil)

	// Mock adding other members
	for _, memberID := range memberIDs {
		mockGroupRepo.On("AddMember", ctx, mock.Anything, memberID, "MEMBER").Return(nil)
	}

	// Execute
	group, err := svc.Create(ctx, creatorID, groupName, memberIDs)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, groupName, group.Name)
	mockGroupRepo.AssertExpectations(t)
}

func TestGroupService_Create_SkipsDuplicateCreator(t *testing.T) {
	ctx := context.Background()
	mockGroupRepo := new(MockGroupRepo)
	svc := service.NewGroupService(mockGroupRepo)

	creatorID := uuid.New()
	// Include creator in member list (should be skipped)
	memberIDs := []uuid.UUID{creatorID, uuid.New(), uuid.New()}
	groupName := "Test Group"

	mockGroupRepo.On("Create", ctx, mock.Anything).Return(nil)
	mockGroupRepo.On("AddMember", ctx, mock.Anything, creatorID, "ADMIN").Return(nil)

	// Should only be called for the 2 non-creator members
	mockGroupRepo.On("AddMember", ctx, mock.Anything, mock.MatchedBy(func(id uuid.UUID) bool {
		return id != creatorID
	}), "MEMBER").Return(nil).Times(2)

	// Execute
	group, err := svc.Create(ctx, creatorID, groupName, memberIDs)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, group)
	mockGroupRepo.AssertExpectations(t)
}

func TestGroupService_AddMember_Success_AsAdmin(t *testing.T) {
	ctx := context.Background()
	mockGroupRepo := new(MockGroupRepo)
	svc := service.NewGroupService(mockGroupRepo)

	adminID := uuid.New()
	groupID := uuid.New()
	newMemberID := uuid.New()

	// Mock: Get members returns admin
	members := []models.GroupMember{
		{GroupID: groupID, UserID: adminID, Role: "ADMIN"},
		{GroupID: groupID, UserID: uuid.New(), Role: "MEMBER"},
	}
	mockGroupRepo.On("GetMembers", ctx, groupID).Return(members, nil)

	// Mock: Check if new member already exists
	mockGroupRepo.On("IsMember", ctx, groupID, newMemberID).Return(false, nil)

	// Mock: Add member
	mockGroupRepo.On("AddMember", ctx, groupID, newMemberID, "MEMBER").Return(nil)

	// Execute
	err := svc.AddMember(ctx, adminID, groupID, newMemberID)

	// Assert
	assert.NoError(t, err)
	mockGroupRepo.AssertExpectations(t)
}

func TestGroupService_AddMember_FailsForNonAdmin(t *testing.T) {
	ctx := context.Background()
	mockGroupRepo := new(MockGroupRepo)
	svc := service.NewGroupService(mockGroupRepo)

	regularUserID := uuid.New()
	groupID := uuid.New()
	newMemberID := uuid.New()

	// Mock: Get members returns only regular members (no admin)
	members := []models.GroupMember{
		{GroupID: groupID, UserID: regularUserID, Role: "MEMBER"},
		{GroupID: groupID, UserID: uuid.New(), Role: "MEMBER"},
	}
	mockGroupRepo.On("GetMembers", ctx, groupID).Return(members, nil)

	// Execute
	err := svc.AddMember(ctx, regularUserID, groupID, newMemberID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "only admins can add members", err.Error())
	mockGroupRepo.AssertExpectations(t)
}

func TestGroupService_AddMember_FailsIfAlreadyMember(t *testing.T) {
	ctx := context.Background()
	mockGroupRepo := new(MockGroupRepo)
	svc := service.NewGroupService(mockGroupRepo)

	adminID := uuid.New()
	groupID := uuid.New()
	existingMemberID := uuid.New()

	// Mock: Get members
	members := []models.GroupMember{
		{GroupID: groupID, UserID: adminID, Role: "ADMIN"},
		{GroupID: groupID, UserID: existingMemberID, Role: "MEMBER"},
	}
	mockGroupRepo.On("GetMembers", ctx, groupID).Return(members, nil)

	// Mock: Member already exists
	mockGroupRepo.On("IsMember", ctx, groupID, existingMemberID).Return(true, nil)

	// Execute
	err := svc.AddMember(ctx, adminID, groupID, existingMemberID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "user is already a member", err.Error())
	mockGroupRepo.AssertExpectations(t)
}

func TestGroupService_RemoveMember_NotYetImplemented(t *testing.T) {
	ctx := context.Background()
	mockGroupRepo := new(MockGroupRepo)
	svc := service.NewGroupService(mockGroupRepo)

	adminID := uuid.New()
	groupID := uuid.New()
	memberID := uuid.New()

	// Mock: Get members
	members := []models.GroupMember{
		{GroupID: groupID, UserID: adminID, Role: "ADMIN"},
		{GroupID: groupID, UserID: memberID, Role: "MEMBER"},
	}
	mockGroupRepo.On("GetMembers", ctx, groupID).Return(members, nil)

	// Execute
	err := svc.RemoveMember(ctx, adminID, groupID, memberID)

	// Assert - should return not implemented error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
}
