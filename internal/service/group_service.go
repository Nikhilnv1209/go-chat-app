package service

import (
	"errors"

	"chat-app/internal/models"
	"chat-app/internal/repository"

	"github.com/google/uuid"
)

type groupService struct {
	groupRepo repository.GroupRepository
}

func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &groupService{
		groupRepo: groupRepo,
	}
}

// Create creates a new group with the given creator as admin and adds the specified members
func (s *groupService) Create(creatorID uuid.UUID, name string, memberIDs []uuid.UUID) (*models.Group, error) {
	// 1. Create the group
	group := &models.Group{
		Name: name,
	}
	if err := s.groupRepo.Create(group); err != nil {
		return nil, err
	}

	// 2. Add creator as ADMIN
	if err := s.groupRepo.AddMember(group.ID, creatorID, "ADMIN"); err != nil {
		return nil, err
	}

	// 3. Add other members
	for _, memberID := range memberIDs {
		// Skip if memberID is the creator (already added)
		if memberID == creatorID {
			continue
		}
		if err := s.groupRepo.AddMember(group.ID, memberID, "MEMBER"); err != nil {
			// Log the error but continue adding other members
			// In production, you might want to rollback or handle this differently
			continue
		}
	}

	return group, nil
}

// AddMember adds a new member to the group (only admins can do this)
func (s *groupService) AddMember(adminID, groupID, newMemberID uuid.UUID) error {
	// 1. Check if the requester is a member and has admin role
	members, err := s.groupRepo.GetMembers(groupID)
	if err != nil {
		return err
	}

	isAdmin := false
	for _, member := range members {
		if member.UserID == adminID && member.Role == "ADMIN" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return errors.New("only admins can add members")
	}

	// 2. Check if the new member is already in the group
	isMember, err := s.groupRepo.IsMember(groupID, newMemberID)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("user is already a member")
	}

	// 3. Add the new member
	return s.groupRepo.AddMember(groupID, newMemberID, "MEMBER")
}

// RemoveMember removes a member from the group (only admins can do this)
func (s *groupService) RemoveMember(adminID, groupID, memberID uuid.UUID) error {
	// 1. Check if the requester is an admin
	members, err := s.groupRepo.GetMembers(groupID)
	if err != nil {
		return err
	}

	isAdmin := false
	for _, member := range members {
		if member.UserID == adminID && member.Role == "ADMIN" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return errors.New("only admins can remove members")
	}

	// 2. Prevent removing yourself if you're the only admin
	if adminID == memberID {
		adminCount := 0
		for _, member := range members {
			if member.Role == "ADMIN" {
				adminCount++
			}
		}
		if adminCount == 1 {
			return errors.New("cannot remove the last admin")
		}
	}

	// 3. Remove the member (for MVP, we'll implement this as a soft delete or skip it)
	// In a real implementation, you'd add a RemoveMember method to the repository
	// For now, we'll return an error indicating it's not yet implemented
	return errors.New("remove member functionality not yet implemented")
}
