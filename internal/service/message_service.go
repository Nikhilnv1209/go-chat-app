package service

import (
	"encoding/json"
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
	msgRepo  repository.MessageRepository
	convRepo repository.ConversationRepository
	hub      Hub
}

func NewMessageService(
	msgRepo repository.MessageRepository,
	convRepo repository.ConversationRepository,
	hub Hub,
) MessageService {
	return &messageService{
		msgRepo:  msgRepo,
		convRepo: convRepo,
		hub:      hub,
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
	// TODO: Implement in F04
	return nil, nil
}

func (s *messageService) GetHistory(userID, targetID uuid.UUID, convType string, limit, beforeID int) ([]models.Message, error) {
	return s.msgRepo.FindByConversation(userID, targetID, convType, limit, beforeID)
}

func (s *messageService) MarkAsRead(userID uuid.UUID, messageIDs []uuid.UUID) error {
	// Placeholder: In a real app we'd mark specific messages.
	// For MVP, we likely just rely on checking the conversation unread count via specific endpoint
	return nil
}
