package service_test

import (
	"chat-app/internal/models"
	"chat-app/internal/service"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMessageRepo
type MockMessageRepo struct {
	mock.Mock
}

func (m *MockMessageRepo) Create(msg *models.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}
func (m *MockMessageRepo) FindByConversation(userID, targetID uuid.UUID, msgType string, limit, beforeID int) ([]models.Message, error) {
	args := m.Called(userID, targetID, msgType, limit, beforeID)
	return args.Get(0).([]models.Message), args.Error(1)
}

// MockConversationRepo
type MockConversationRepo struct {
	mock.Mock
}

func (m *MockConversationRepo) Upsert(conv *models.Conversation) error {
	args := m.Called(conv)
	return args.Error(0)
}
func (m *MockConversationRepo) FindByUser(userID uuid.UUID) ([]models.Conversation, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Conversation), args.Error(1)
}
func (m *MockConversationRepo) IncrementUnread(userID uuid.UUID, convType string, targetID uuid.UUID) error {
	args := m.Called(userID, convType, targetID)
	return args.Error(0)
}
func (m *MockConversationRepo) ResetUnread(userID uuid.UUID, convType string, targetID uuid.UUID) error {
	args := m.Called(userID, convType, targetID)
	return args.Error(0)
}

// MockHub
type MockHub struct {
	mock.Mock
}

func (m *MockHub) SendToUser(userID uuid.UUID, message []byte) {
	m.Called(userID, message)
}

func TestSendDirectMessage(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockHub := new(MockHub)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockHub)

	senderID := uuid.New()
	receiverID := uuid.New()
	content := "Hello"

	// Expectations
	// 1. Create Message
	mockMsgRepo.On("Create", mock.MatchedBy(func(msg *models.Message) bool {
		return msg.SenderID == senderID && *msg.ReceiverID == receiverID && msg.Content == content
	})).Return(nil)

	// 2. Upsert Conversation for Sender
	mockConvRepo.On("Upsert", mock.MatchedBy(func(conv *models.Conversation) bool {
		return conv.UserID == senderID && conv.TargetID == receiverID
	})).Return(nil)

	// 3. Increment Unread for Receiver
	mockConvRepo.On("IncrementUnread", receiverID, "private", senderID).Return(nil)

	// 4. Send to Hub using MatchedBy to ignore dynamic JSON
	mockHub.On("SendToUser", receiverID, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)
		return msg["type"] == "new_message"
	})).Return()

	// Execute
	msg, err := svc.SendDirectMessage(senderID, receiverID, content)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, content, msg.Content)

	mockMsgRepo.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
	mockHub.AssertExpectations(t)
}

func TestGetHistory_Conversation(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockHub := new(MockHub)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockHub)

	userID := uuid.New()
	targetID := uuid.New()
	limit := 10

	// Create 10 mock messages
	mockMessages := make([]models.Message, 10)
	for i := 0; i < 10; i++ {
		mockMessages[i] = models.Message{
			BaseModel:  models.BaseModel{ID: uuid.New(), CreatedAt: time.Now()},
			SenderID:   userID,
			ReceiverID: &targetID,
			Content:    "Message " + string(rune(i)),
			MsgType:    "private",
		}
	}

	mockMsgRepo.On("FindByConversation", userID, targetID, "private", limit, 0).Return(mockMessages, nil)

	history, err := svc.GetHistory(userID, targetID, "private", limit, 0)
	assert.NoError(t, err)
	assert.Len(t, history, 10)
	assert.Equal(t, "Message \x00", history[0].Content)

	mockMsgRepo.AssertExpectations(t)
}

func TestLongConversationFlow(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockHub := new(MockHub)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockHub)

	user1 := uuid.New()
	user2 := uuid.New()

	// Simulate 10 messages back and forth
	for i := 0; i < 10; i++ {
		sender := user1
		receiver := user2
		if i%2 == 1 {
			sender = user2
			receiver = user1
		}
		content := "Message " + string(rune(i))

		// Expectations for each message
		mockMsgRepo.On("Create", mock.MatchedBy(func(msg *models.Message) bool {
			return msg.SenderID == sender && *msg.ReceiverID == receiver && msg.Content == content
		})).Return(nil).Once()

		mockConvRepo.On("Upsert", mock.MatchedBy(func(conv *models.Conversation) bool {
			return conv.UserID == sender && conv.TargetID == receiver
		})).Return(nil).Once()

		mockConvRepo.On("IncrementUnread", receiver, "private", sender).Return(nil).Once()

		mockHub.On("SendToUser", receiver, mock.Anything).Return().Once()

		// Execute
		_, err := svc.SendDirectMessage(sender, receiver, content)
		assert.NoError(t, err)
	}

	mockMsgRepo.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
	mockHub.AssertExpectations(t)
}
