package service

import (
	"context"
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
	msgRepo     repository.MessageRepository
	convRepo    repository.ConversationRepository
	groupRepo   repository.GroupRepository
	receiptRepo repository.MessageReceiptRepository
	userRepo    repository.UserRepository
	hub         Hub
}

func NewMessageService(
	msgRepo repository.MessageRepository,
	convRepo repository.ConversationRepository,
	groupRepo repository.GroupRepository,
	receiptRepo repository.MessageReceiptRepository,
	userRepo repository.UserRepository,
	hub Hub,
) MessageService {
	return &messageService{
		msgRepo:     msgRepo,
		convRepo:    convRepo,
		groupRepo:   groupRepo,
		receiptRepo: receiptRepo,
		userRepo:    userRepo,
		hub:         hub,
	}
}

func (s *messageService) SendDirectMessage(ctx context.Context, senderID, receiverID uuid.UUID, content string) (*models.Message, error) {
	// 1. Create Message
	msg := &models.Message{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
		},
		SenderID:   senderID,
		ReceiverID: &receiverID,
		Content:    content,
		MsgType:    "DM",
	}

	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	// 1.5 Populate Sender info for response/broadcast
	sender, err := s.userRepo.FindByID(ctx, senderID)
	if err == nil {
		msg.Sender = *sender
	}

	// 2. Create Receipt with SENT status [F06]
	receipt := &models.MessageReceipt{
		MessageID: msg.ID,
		UserID:    receiverID,
		Status:    "SENT",
	}
	if err := s.receiptRepo.Create(ctx, receipt); err != nil {
		// Log error but don't fail the message send
		// In production, use proper logging
	}

	// 3. Update Conversations (Sender)
	// Upsert Sender's conversation with Receiver
	s.convRepo.Upsert(ctx, &models.Conversation{
		UserID:        senderID,
		Type:          "DM",
		TargetID:      receiverID,
		LastMessage:   content,
		LastMessageAt: msg.CreatedAt,
		UnreadCount:   0, // Sender doesn't have unread
	})

	// 4. Update Conversations (Receiver)
	// Increment Receiver's unread count for conversation with Sender
	s.convRepo.IncrementUnread(ctx, receiverID, "DM", senderID, content)

	// 5. Real-time Delivery via WebSocket
	payload, _ := json.Marshal(map[string]interface{}{
		"type":    "new_message",
		"payload": msg,
	})

	// Broadcast to receiver's devices
	s.hub.SendToUser(receiverID, payload)

	// Broadcast to sender's other devices for multi-device sync
	s.hub.SendToUser(senderID, payload)

	return msg, nil
}

func (s *messageService) SendGroupMessage(ctx context.Context, senderID, groupID uuid.UUID, content string) (*models.Message, error) {
	// 1. Verify sender is a member of the group
	isMember, err := s.groupRepo.IsMember(ctx, groupID, senderID)
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
		MsgType:  "GROUP",
	}

	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	// 2.5 Populate Sender info for response/broadcast
	sender, err := s.userRepo.FindByID(ctx, senderID)
	if err == nil {
		msg.Sender = *sender
	}

	// 3. Create Receipts for all members except sender [F06] - Batch Optimized
	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var receipts []*models.MessageReceipt
	for _, member := range members {
		if member.UserID != senderID {
			receipt := &models.MessageReceipt{
				BaseModel: models.BaseModel{
					CreatedAt: time.Now(),
				},
				MessageID: msg.ID,
				UserID:    member.UserID,
				Status:    "SENT",
				UpdatedAt: time.Now(),
			}
			receipts = append(receipts, receipt)
		}
	}

	if len(receipts) > 0 {
		if err := s.receiptRepo.CreateBatch(ctx, receipts); err != nil {
			// Log error but continue? Or fail?
			// Ideally fail if data integrity is key, but for chat performant delivery
			// we might log. For now, strict: return error.
			return nil, err
		}
	}

	// 4. For each member (except sender), update conversation and broadcast
	for _, member := range members {

		if member.UserID == senderID {
			// Update sender's conversation without incrementing unread
			s.convRepo.Upsert(ctx, &models.Conversation{
				UserID:        senderID,
				Type:          "GROUP",
				TargetID:      groupID,
				LastMessage:   content,
				LastMessageAt: msg.CreatedAt,
				UnreadCount:   0,
			})
			continue
		}

		// Update receiver's conversation and increment unread
		s.convRepo.IncrementUnread(ctx, member.UserID, "GROUP", groupID, content)

		// Real-time delivery via WebSocket
		payload, _ := json.Marshal(map[string]interface{}{
			"type":    "new_message",
			"payload": msg,
		})
		s.hub.SendToUser(member.UserID, payload)
	}

	// Broadcast to sender's devices for multi-device sync
	senderPayload, _ := json.Marshal(map[string]interface{}{
		"type":    "new_message",
		"payload": msg,
	})
	s.hub.SendToUser(senderID, senderPayload)

	return msg, nil
}

// Custom error for group messaging
var ErrNotGroupMember = errors.New("sender is not a member of the group")

func (s *messageService) GetHistory(ctx context.Context, userID, targetID uuid.UUID, convType string, limit int, beforeID *uuid.UUID) ([]models.Message, error) {
	return s.msgRepo.FindByConversation(ctx, userID, targetID, convType, limit, beforeID)
}

func (s *messageService) MarkAsRead(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) error {
	return s.updateReceiptStatus(ctx, userID, messageIDs, "READ")
}

func (s *messageService) MarkAsDelivered(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) error {
	return s.updateReceiptStatus(ctx, userID, messageIDs, "DELIVERED")
}

func (s *messageService) updateReceiptStatus(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID, status string) error {
	for _, msgID := range messageIDs {
		// 1. Update Receipt Status
		if err := s.receiptRepo.UpdateStatus(ctx, msgID, userID, status); err != nil {
			// Skip or log error, but continue for other messages
			continue
		}

		// 2. Find Message to identify Sender
		msg, err := s.msgRepo.FindByID(ctx, msgID)
		if err != nil {
			continue
		}

		// 3. Broadcast receipt_update to Sender [F06]
		// Don't broadcast to self if it's a note to self (unlikely)
		if msg.SenderID != userID {
			payload, _ := json.Marshal(map[string]interface{}{
				"type": "receipt_update",
				"payload": map[string]interface{}{
					"message_id": msgID,
					"user_id":    userID, // Who read/received it
					"status":     status,
					"updated_at": time.Now(),
				},
			})
			s.hub.SendToUser(msg.SenderID, payload)
		}
	}
	return nil
}

func (s *messageService) GetMessageReceipts(ctx context.Context, userID, messageID uuid.UUID) ([]models.MessageReceipt, error) {
	// 1. Fetch Message to verify access
	msg, err := s.msgRepo.FindByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	// 2. Validate Access
	hasAccess := false
	if msg.SenderID == userID {
		hasAccess = true
	} else if msg.ReceiverID != nil && *msg.ReceiverID == userID {
		hasAccess = true
	} else if msg.GroupID != nil {
		isMember, err := s.groupRepo.IsMember(ctx, *msg.GroupID, userID)
		if err == nil && isMember {
			hasAccess = true
		}
	}

	if !hasAccess {
		return nil, errors.New("access denied")
	}

	// 3. Fetch Receipts
	return s.receiptRepo.FindByMessageID(ctx, messageID)
}

func (s *messageService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *messageService) BroadcastTypingIndicator(ctx context.Context, userID uuid.UUID, username, convType string, targetID uuid.UUID, isTyping bool) error {
	// Build event type
	eventType := "user_typing"
	if !isTyping {
		eventType = "user_stopped_typing"
	}

	// Build base payload
	payloadData := map[string]interface{}{
		"user_id":           userID,
		"conversation_type": convType,
		"target_id":         targetID,
	}

	// Add username only for typing_start events
	if isTyping {
		payloadData["username"] = username
	}

	switch convType {
	case "DM":
		// Send to single user
		// Prevent broadcast to self
		if targetID != userID {
			payload, _ := json.Marshal(map[string]interface{}{
				"type":    eventType,
				"payload": payloadData,
			})
			s.hub.SendToUser(targetID, payload)
		}
	case "GROUP":
		// Verify sender is a member
		isMember, err := s.groupRepo.IsMember(ctx, targetID, userID)
		if err != nil {
			return err
		}
		if !isMember {
			return ErrNotGroupMember
		}

		// Get all group members
		members, err := s.groupRepo.GetMembers(ctx, targetID)
		if err != nil {
			return err
		}

		// Broadcast to all members except sender
		payload, _ := json.Marshal(map[string]interface{}{
			"type":    eventType,
			"payload": payloadData,
		})

		for _, member := range members {
			if member.UserID != userID {
				s.hub.SendToUser(member.UserID, payload)
			}
		}
	}

	return nil
}
