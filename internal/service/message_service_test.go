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

func (m *MockMessageRepo) FindByID(id uuid.UUID) (*models.Message, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
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

// MockGroupRepo
type MockGroupRepo struct {
	mock.Mock
}

func (m *MockGroupRepo) Create(group *models.Group) error {
	args := m.Called(group)
	return args.Error(0)
}

func (m *MockGroupRepo) FindByID(id uuid.UUID) (*models.Group, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *MockGroupRepo) GetMembers(groupID uuid.UUID) ([]models.GroupMember, error) {
	args := m.Called(groupID)
	return args.Get(0).([]models.GroupMember), args.Error(1)
}

func (m *MockGroupRepo) IsMember(groupID, userID uuid.UUID) (bool, error) {
	args := m.Called(groupID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockGroupRepo) AddMember(groupID, userID uuid.UUID, role string) error {
	args := m.Called(groupID, userID, role)
	return args.Error(0)
}

// MockHub
type MockHub struct {
	mock.Mock
}

func (m *MockHub) SendToUser(userID uuid.UUID, message []byte) {
	m.Called(userID, message)
}

// MockMessageReceiptRepo [F06]
type MockMessageReceiptRepo struct {
	mock.Mock
}

func (m *MockMessageReceiptRepo) Create(receipt *models.MessageReceipt) error {
	args := m.Called(receipt)
	return args.Error(0)
}

func (m *MockMessageReceiptRepo) CreateBatch(receipts []*models.MessageReceipt) error {
	args := m.Called(receipts)
	return args.Error(0)
}

func (m *MockMessageReceiptRepo) UpdateStatus(messageID, userID uuid.UUID, status string) error {
	args := m.Called(messageID, userID, status)
	return args.Error(0)
}

func (m *MockMessageReceiptRepo) FindByMessageID(messageID uuid.UUID) ([]models.MessageReceipt, error) {
	args := m.Called(messageID)
	return args.Get(0).([]models.MessageReceipt), args.Error(1)
}

func (m *MockMessageReceiptRepo) FindUnreadCount(userID uuid.UUID) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func TestSendDirectMessage(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

	senderID := uuid.New()
	receiverID := uuid.New()
	content := "Hello"

	// Expectations
	// 1. Create Message
	mockMsgRepo.On("Create", mock.MatchedBy(func(msg *models.Message) bool {
		return msg.SenderID == senderID && *msg.ReceiverID == receiverID && msg.Content == content
	})).Return(nil)

	// 2. Create Receipt [F06]
	mockReceiptRepo.On("Create", mock.MatchedBy(func(receipt *models.MessageReceipt) bool {
		return receipt.UserID == receiverID && receipt.Status == "SENT"
	})).Return(nil)

	// 3. Upsert Conversation for Sender
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
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

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
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

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

		// Receipt creation [F06]
		mockReceiptRepo.On("Create", mock.MatchedBy(func(receipt *models.MessageReceipt) bool {
			return receipt.UserID == receiver && receipt.Status == "SENT"
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

// ===== GROUP MESSAGING TESTS =====

func TestSendGroupMessage_Success(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

	senderID := uuid.New()
	groupID := uuid.New()
	member1 := uuid.New()
	member2 := uuid.New()
	content := "Hello Group!"

	// Mock: Sender is a member
	mockGroupRepo.On("IsMember", groupID, senderID).Return(true, nil)

	// Mock: Create message
	mockMsgRepo.On("Create", mock.MatchedBy(func(msg *models.Message) bool {
		return msg.SenderID == senderID && msg.GroupID != nil && *msg.GroupID == groupID && msg.Content == content
	})).Return(nil)

	// Mock: Create Receipts [F06]
	// Mock: Create Receipts [F06] - Batch
	mockReceiptRepo.On("CreateBatch", mock.MatchedBy(func(receipts []*models.MessageReceipt) bool {
		return len(receipts) == 2 && receipts[0].Status == "SENT"
	})).Return(nil).Once()

	// Mock: Get group members
	members := []models.GroupMember{
		{GroupID: groupID, UserID: senderID, Role: "ADMIN"},
		{GroupID: groupID, UserID: member1, Role: "MEMBER"},
		{GroupID: groupID, UserID: member2, Role: "MEMBER"},
	}
	mockGroupRepo.On("GetMembers", groupID).Return(members, nil)

	// Mock: Upsert sender's conversation
	mockConvRepo.On("Upsert", mock.MatchedBy(func(conv *models.Conversation) bool {
		return conv.UserID == senderID && conv.TargetID == groupID && conv.Type == "group"
	})).Return(nil)

	// Mock: Increment unread for other members
	mockConvRepo.On("IncrementUnread", member1, "group", groupID).Return(nil)
	mockConvRepo.On("IncrementUnread", member2, "group", groupID).Return(nil)

	// Mock: Send to hub for other members
	mockHub.On("SendToUser", member1, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)
		return msg["type"] == "new_message"
	})).Return()

	mockHub.On("SendToUser", member2, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)
		return msg["type"] == "new_message"
	})).Return()

	// Execute
	msg, err := svc.SendGroupMessage(senderID, groupID, content)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, content, msg.Content)
	assert.Equal(t, groupID, *msg.GroupID)

	mockMsgRepo.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
	mockGroupRepo.AssertExpectations(t)
	mockHub.AssertExpectations(t)
}

func TestSendGroupMessage_FailsForNonMember(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

	nonMemberID := uuid.New()
	groupID := uuid.New()
	content := "I shouldn't be able to send this"

	// Mock: Sender is NOT a member
	mockGroupRepo.On("IsMember", groupID, nonMemberID).Return(false, nil)

	// Execute
	msg, err := svc.SendGroupMessage(nonMemberID, groupID, content)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, msg)
	assert.Equal(t, "sender is not a member of the group", err.Error())

	mockGroupRepo.AssertExpectations(t)
	// Message should NOT be created
	mockMsgRepo.AssertNotCalled(t, "Create")
}

func TestSendGroupMessage_BroadcastsToAllMembers(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

	senderID := uuid.New()
	groupID := uuid.New()
	// Create 5 members total (sender + 4 others)
	members := []models.GroupMember{
		{GroupID: groupID, UserID: senderID, Role: "ADMIN"},
	}
	otherMemberIDs := []uuid.UUID{}
	for i := 0; i < 4; i++ {
		memberID := uuid.New()
		members = append(members, models.GroupMember{
			GroupID: groupID,
			UserID:  memberID,
			Role:    "MEMBER",
		})
		otherMemberIDs = append(otherMemberIDs, memberID)
	}

	// Mock setup
	mockGroupRepo.On("IsMember", groupID, senderID).Return(true, nil)
	mockMsgRepo.On("Create", mock.Anything).Return(nil)
	// Mock: Create Receipts [F06]
	// Mock: Create Receipts [F06] - Batch
	mockReceiptRepo.On("CreateBatch", mock.MatchedBy(func(receipts []*models.MessageReceipt) bool {
		return len(receipts) == 4
	})).Return(nil).Once()
	mockGroupRepo.On("GetMembers", groupID).Return(members, nil)
	mockConvRepo.On("Upsert", mock.Anything).Return(nil)

	// Mock for each other member
	for _, memberID := range otherMemberIDs {
		mockConvRepo.On("IncrementUnread", memberID, "group", groupID).Return(nil)
		mockHub.On("SendToUser", memberID, mock.Anything).Return()
	}

	// Execute
	msg, err := svc.SendGroupMessage(senderID, groupID, "Test broadcast")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, msg)

	// Verify each other member received the message
	for _, memberID := range otherMemberIDs {
		mockHub.AssertCalled(t, "SendToUser", memberID, mock.Anything)
	}

	mockGroupRepo.AssertExpectations(t)
}

func TestSendGroupMessage_UpdatesConversationForAllMembers(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

	senderID := uuid.New()
	member1 := uuid.New()
	member2 := uuid.New()
	groupID := uuid.New()

	members := []models.GroupMember{
		{GroupID: groupID, UserID: senderID, Role: "MEMBER"},
		{GroupID: groupID, UserID: member1, Role: "MEMBER"},
		{GroupID: groupID, UserID: member2, Role: "MEMBER"},
	}

	mockGroupRepo.On("IsMember", groupID, senderID).Return(true, nil)
	mockMsgRepo.On("Create", mock.Anything).Return(nil)
	// Mock: Create Receipts [F06]
	// Mock: Create Receipts [F06] - Batch
	mockReceiptRepo.On("CreateBatch", mock.MatchedBy(func(receipts []*models.MessageReceipt) bool {
		return len(receipts) == 2
	})).Return(nil).Once()
	mockGroupRepo.On("GetMembers", groupID).Return(members, nil)

	// Expect conversation upsert for sender (unread = 0)
	mockConvRepo.On("Upsert", mock.MatchedBy(func(conv *models.Conversation) bool {
		return conv.UserID == senderID && conv.UnreadCount == 0
	})).Return(nil)

	// Expect increment unread for other members
	mockConvRepo.On("IncrementUnread", member1, "group", groupID).Return(nil)
	mockConvRepo.On("IncrementUnread", member2, "group", groupID).Return(nil)

	mockHub.On("SendToUser", mock.Anything, mock.Anything).Return()

	// Execute
	_, err := svc.SendGroupMessage(senderID, groupID, "Update conversations")

	// Assert
	assert.NoError(t, err)

	// Verify conversations were updated
	mockConvRepo.AssertCalled(t, "Upsert", mock.MatchedBy(func(conv *models.Conversation) bool {
		return conv.UserID == senderID
	}))
	mockConvRepo.AssertCalled(t, "IncrementUnread", member1, "group", groupID)
	mockConvRepo.AssertCalled(t, "IncrementUnread", member2, "group", groupID)
}

func TestSendGroupMessage_SenderDoesNotReceiveOwnMessage(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

	senderID := uuid.New()
	member1 := uuid.New()
	groupID := uuid.New()

	members := []models.GroupMember{
		{GroupID: groupID, UserID: senderID, Role: "ADMIN"},
		{GroupID: groupID, UserID: member1, Role: "MEMBER"},
	}

	mockGroupRepo.On("IsMember", groupID, senderID).Return(true, nil)
	mockMsgRepo.On("Create", mock.Anything).Return(nil)
	// Mock: Create Receipts [F06]
	// Mock: Create Receipts [F06] - Batch
	mockReceiptRepo.On("CreateBatch", mock.MatchedBy(func(receipts []*models.MessageReceipt) bool {
		return len(receipts) == 1
	})).Return(nil).Once()
	mockGroupRepo.On("GetMembers", groupID).Return(members, nil)
	mockConvRepo.On("Upsert", mock.Anything).Return(nil)
	mockConvRepo.On("IncrementUnread", member1, "group", groupID).Return(nil)
	mockHub.On("SendToUser", member1, mock.Anything).Return()

	// Execute
	_, err := svc.SendGroupMessage(senderID, groupID, "Own message test")

	// Assert
	assert.NoError(t, err)

	// Sender should NOT receive the message via hub
	mockHub.AssertNotCalled(t, "SendToUser", senderID, mock.Anything)
	// But member1 should receive it
	mockHub.AssertCalled(t, "SendToUser", member1, mock.Anything)
}

func TestGetMessageReceipts_Success_Sender(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)
	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

	userID := uuid.New()
	msgID := uuid.New()

	msg := &models.Message{
		BaseModel: models.BaseModel{ID: msgID},
		SenderID:  userID,
	}

	mockMsgRepo.On("FindByID", msgID).Return(msg, nil)

	mockReceipts := []models.MessageReceipt{{Status: "READ"}}
	mockReceiptRepo.On("FindByMessageID", msgID).Return(mockReceipts, nil)

	res, err := svc.GetMessageReceipts(userID, msgID)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestGetMessageReceipts_Forbidden(t *testing.T) {
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)
	mockReceiptRepo := new(MockMessageReceiptRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockHub)

	userID := uuid.New()
	otherUser := uuid.New()
	msgID := uuid.New()

	msg := &models.Message{
		BaseModel: models.BaseModel{ID: msgID},
		SenderID:  otherUser, // Not sender
		// No ReceiverID, No GroupID -> generic forbidden
	}

	mockMsgRepo.On("FindByID", msgID).Return(msg, nil)

	_, err := svc.GetMessageReceipts(userID, msgID)
	assert.Error(t, err)
	assert.Equal(t, "access denied", err.Error())
}
