package service

import (
	"encoding/json"
	"errors"
	"time"

	"chat-app/internal/models"
	"chat-app/internal/repository"

	"github.com/google/uuid"
)

// Hub defines the methods needed by the MessageService to broadcast messages.
// This decouples the service package from the websocket package.
type Hub interface {
	SendToUser(userID uuid.UUID, message []byte)
}

type messageService struct {
	msgRepo   repository.MessageRepository
	convRepo  repository.ConversationRepository
	groupRepo repository.GroupRepository
	hub       Hub
}

func NewMessageService(
	msgRepo repository.MessageRepository,
	convRepo repository.ConversationRepository,
	groupRepo repository.GroupRepository,
	hub Hub,
) MessageService {
	return &messageService{
		msgRepo:   msgRepo,
		convRepo:  convRepo,
		groupRepo: groupRepo,
		hub:       hub,
	}
}

func (s *messageService) SendDirectMessage(senderID, receiverID uuid.UUID, content string) (*models.Message, error) {
	// 1. Create Message
	msg := &models.Message{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
		},
		SenderID:   senderID,
		ReceiverID: &receiverID,
		Content:    content,
		MsgType:    "private",
	}

	if err := s.msgRepo.Create(msg); err != nil {
		return nil, err
	}

	// 2. Update Conversations (Sender)
	// Upsert Sender's conversation with Receiver
	s.convRepo.Upsert(&models.Conversation{
		UserID:        senderID,
		Type:          "private",
		TargetID:      receiverID,
		LastMessageAt: msg.CreatedAt,
		UnreadCount:   0, // Sender doesn't have unread
	})

	// 3. Update Conversations (Receiver)
	// Increment Receiver's unread count for conversation with Sender
	s.convRepo.IncrementUnread(receiverID, "private", senderID)

	// 4. Real-time Delivery via WebSocket
	payload, _ := json.Marshal(map[string]interface{}{
		"type":    "new_message",
		"payload": msg,
	})
	s.hub.SendToUser(receiverID, payload)

	return msg, nil
}

func (s *messageService) SendGroupMessage(senderID, groupID uuid.UUID, content string) (*models.Message, error) {
	// 1. Verify sender is a member of the group
	isMember, err := s.groupRepo.IsMember(groupID, senderID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotGroupMember
	}

	// 2. Create Message
	msg := &models.Message{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
		},
		SenderID: senderID,
		GroupID:  &groupID,
		Content:  content,
		MsgType:  "group",
	}

	if err := s.msgRepo.Create(msg); err != nil {
		return nil, err
	}

	// 3. Get all group members
	members, err := s.groupRepo.GetMembers(groupID)
	if err != nil {
		return nil, err
	}

	// 4. For each member (except sender), update conversation and broadcast
	for _, member := range members {
		if member.UserID == senderID {
			// Update sender's conversation without incrementing unread
			s.convRepo.Upsert(&models.Conversation{
				UserID:        senderID,
				Type:          "group",
				TargetID:      groupID,
				LastMessageAt: msg.CreatedAt,
				UnreadCount:   0,
			})
			continue
		}

		// Update receiver's conversation and increment unread
		s.convRepo.IncrementUnread(member.UserID, "group", groupID)

		// Real-time delivery via WebSocket
		payload, _ := json.Marshal(map[string]interface{}{
			"type":    "new_message",
			"payload": msg,
		})
		s.hub.SendToUser(member.UserID, payload)
	}

	return msg, nil
}

// Custom error for group messaging
var ErrNotGroupMember = errors.New("sender is not a member of the group")

func (s *messageService) GetHistory(userID, targetID uuid.UUID, convType string, limit, beforeID int) ([]models.Message, error) {
	return s.msgRepo.FindByConversation(userID, targetID, convType, limit, beforeID)
}

func (s *messageService) MarkAsRead(userID uuid.UUID, messageIDs []uuid.UUID) error {
	// Placeholder: In a real app we'd mark specific messages.
	// For MVP, we likely just rely on checking the conversation unread count via specific endpoint
	return nil
}
