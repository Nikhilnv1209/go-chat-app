package service_test

import (
	"context"
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

func (m *MockMessageRepo) Create(ctx context.Context, msg *models.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}
func (m *MockMessageRepo) FindByConversation(ctx context.Context, userID, targetID uuid.UUID, msgType string, limit int, beforeID *uuid.UUID) ([]models.Message, error) {
	args := m.Called(ctx, userID, targetID, msgType, limit, beforeID)
	return args.Get(0).([]models.Message), args.Error(1)
}

func (m *MockMessageRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

// MockConversationRepo
type MockConversationRepo struct {
	mock.Mock
}

func (m *MockConversationRepo) Upsert(ctx context.Context, conv *models.Conversation) error {
	args := m.Called(ctx, conv)
	return args.Error(0)
}
func (m *MockConversationRepo) FindByUser(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Conversation), args.Error(1)
}
func (m *MockConversationRepo) IncrementUnread(ctx context.Context, userID uuid.UUID, convType string, targetID uuid.UUID, lastMessage string) error {
	args := m.Called(ctx, userID, convType, targetID, lastMessage)
	return args.Error(0)
}
func (m *MockConversationRepo) ResetUnread(ctx context.Context, userID uuid.UUID, convType string, targetID uuid.UUID) error {
	args := m.Called(ctx, userID, convType, targetID)
	return args.Error(0)
}

func (m *MockConversationRepo) FindContactsOfUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

// MockGroupRepo
type MockGroupRepo struct {
	mock.Mock
}

func (m *MockGroupRepo) Create(ctx context.Context, group *models.Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockGroupRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *MockGroupRepo) GetMembers(ctx context.Context, groupID uuid.UUID) ([]models.GroupMember, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]models.GroupMember), args.Error(1)
}

func (m *MockGroupRepo) IsMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, groupID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockGroupRepo) AddMember(ctx context.Context, groupID, userID uuid.UUID, role string) error {
	args := m.Called(ctx, groupID, userID, role)
	return args.Error(0)
}

// MockHub
type MockHub struct {
	mock.Mock
}

func (m *MockHub) SendToUser(userID uuid.UUID, message []byte) {
	m.Called(userID, message)
}

func (m *MockHub) IsUserViewingConversation(convType string, targetID uuid.UUID) bool {
	args := m.Called(convType, targetID)
	return args.Bool(0)
}

// MockUserRepo [F07]
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) UpdateOnlineStatus(ctx context.Context, userID uuid.UUID, isOnline bool, lastSeen time.Time) error {
	args := m.Called(ctx, userID, isOnline, lastSeen)
	return args.Error(0)
}

func (m *MockUserRepo) Search(ctx context.Context, query string, excludeUserID uuid.UUID) ([]models.User, error) {
	args := m.Called(ctx, query, excludeUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

// MockMessageReceiptRepo [F06]
type MockMessageReceiptRepo struct {
	mock.Mock
}

func (m *MockMessageReceiptRepo) Create(ctx context.Context, receipt *models.MessageReceipt) error {
	args := m.Called(ctx, receipt)
	return args.Error(0)
}

func (m *MockMessageReceiptRepo) CreateBatch(ctx context.Context, receipts []*models.MessageReceipt) error {
	args := m.Called(ctx, receipts)
	return args.Error(0)
}

func (m *MockMessageReceiptRepo) UpdateStatus(ctx context.Context, messageID, userID uuid.UUID, status string) error {
	args := m.Called(ctx, messageID, userID, status)
	return args.Error(0)
}

func (m *MockMessageReceiptRepo) FindByMessageID(ctx context.Context, messageID uuid.UUID) ([]models.MessageReceipt, error) {
	args := m.Called(ctx, messageID)
	return args.Get(0).([]models.MessageReceipt), args.Error(1)
}

func (m *MockMessageReceiptRepo) FindUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}


func TestSendDirectMessage(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	senderID := uuid.New()
	receiverID := uuid.New()
	content := "Hello"

	// Expectations
	// 1. Create Message
	mockMsgRepo.On("Create", ctx, mock.MatchedBy(func(msg *models.Message) bool {
		return msg.SenderID == senderID && *msg.ReceiverID == receiverID && msg.Content == content
	})).Return(nil)

	// 1.5 Get sender info for response
	mockUserRepo.On("FindByID", ctx, senderID).Return(&models.User{
		BaseModel: models.BaseModel{ID: senderID},
		Username:  "Sender",
	}, nil)

	// 2. Create Receipt [F06]
	mockReceiptRepo.On("Create", ctx, mock.MatchedBy(func(receipt *models.MessageReceipt) bool {
		return receipt.UserID == receiverID && receipt.Status == "SENT"
	})).Return(nil)

	// 3. Upsert Conversation for Sender
	mockConvRepo.On("Upsert", ctx, mock.MatchedBy(func(conv *models.Conversation) bool {
		return conv.UserID == senderID && conv.TargetID == receiverID
	})).Return(nil)

	// 3.5 Check if receiver is viewing the conversation (returns false by default for this test)
	mockHub.On("IsUserViewingConversation", "DM", senderID).Return(false)

	// 3. Increment Unread for Receiver (now with 5 params including lastMessage)
	mockConvRepo.On("IncrementUnread", ctx, receiverID, "DM", senderID, content).Return(nil)

	// 4. Send to Hub for receiver (B006: also send to sender's other devices)
	mockHub.On("SendToUser", receiverID, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)
		return msg["type"] == "new_message"
	})).Return()

	// B006: Also broadcast to sender's other devices for multi-device sync
	mockHub.On("SendToUser", senderID, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)
		return msg["type"] == "new_message"
	})).Return()

	// Execute
	msg, err := svc.SendDirectMessage(ctx, senderID, receiverID, content)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, content, msg.Content)

	mockMsgRepo.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
	mockHub.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetHistory_Conversation(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

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
			MsgType:    "DM",
		}
	}

	mockMsgRepo.On("FindByConversation", ctx, userID, targetID, "DM", limit, (*uuid.UUID)(nil)).Return(mockMessages, nil)

	history, err := svc.GetHistory(ctx, userID, targetID, "DM", limit, nil)
	assert.NoError(t, err)
	assert.Len(t, history, 10)
	assert.Equal(t, "Message \x00", history[0].Content)

	mockMsgRepo.AssertExpectations(t)
}

func TestLongConversationFlow(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

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
		mockMsgRepo.On("Create", ctx, mock.MatchedBy(func(msg *models.Message) bool {
			return msg.SenderID == sender && *msg.ReceiverID == receiver && msg.Content == content
		})).Return(nil).Once()

		// Get sender info for response
		mockUserRepo.On("FindByID", ctx, sender).Return(&models.User{
			BaseModel: models.BaseModel{ID: sender},
			Username:  "User",
		}, nil).Once()

		// Receipt creation [F06]
		mockReceiptRepo.On("Create", ctx, mock.MatchedBy(func(receipt *models.MessageReceipt) bool {
			return receipt.UserID == receiver && receipt.Status == "SENT"
		})).Return(nil).Once()

		mockConvRepo.On("Upsert", ctx, mock.MatchedBy(func(conv *models.Conversation) bool {
			return conv.UserID == sender && conv.TargetID == receiver
		})).Return(nil).Once()

		// Check if viewing conversation (returns false by default)
		mockHub.On("IsUserViewingConversation", "DM", sender).Return(false).Once()

		mockConvRepo.On("IncrementUnread", ctx, receiver, "DM", sender, content).Return(nil).Once()

		// B006: Send to both receiver and sender (for multi-device sync)
		mockHub.On("SendToUser", receiver, mock.Anything).Return().Once()
		mockHub.On("SendToUser", sender, mock.Anything).Return().Once()

		// Execute
		_, err := svc.SendDirectMessage(ctx, sender, receiver, content)
		assert.NoError(t, err)
	}

	mockMsgRepo.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
	mockHub.AssertExpectations(t)
}

// ===== GROUP MESSAGING TESTS =====

func TestSendGroupMessage_Success(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	senderID := uuid.New()
	groupID := uuid.New()
	member1 := uuid.New()
	member2 := uuid.New()
	content := "Hello Group!"

	// Mock: Sender is a member
	mockGroupRepo.On("IsMember", ctx, groupID, senderID).Return(true, nil)

	// Mock: Create message
	mockMsgRepo.On("Create", ctx, mock.MatchedBy(func(msg *models.Message) bool {
		return msg.SenderID == senderID && msg.GroupID != nil && *msg.GroupID == groupID && msg.Content == content
	})).Return(nil)

	// Mock: Get sender info for response
	mockUserRepo.On("FindByID", ctx, senderID).Return(&models.User{
		BaseModel: models.BaseModel{ID: senderID},
		Username:  "Sender",
	}, nil)

	// Mock: Create Receipts [F06]
	// Mock: Create Receipts [F06] - Batch
	mockReceiptRepo.On("CreateBatch", ctx, mock.MatchedBy(func(receipts []*models.MessageReceipt) bool {
		return len(receipts) == 2 && receipts[0].Status == "SENT"
	})).Return(nil).Once()

	// Mock: Get group members
	members := []models.GroupMember{
		{GroupID: groupID, UserID: senderID, Role: "ADMIN"},
		{GroupID: groupID, UserID: member1, Role: "MEMBER"},
		{GroupID: groupID, UserID: member2, Role: "MEMBER"},
	}
	mockGroupRepo.On("GetMembers", ctx, groupID).Return(members, nil)

	// Mock: Check if viewing conversation (returns false for all members)
	mockHub.On("IsUserViewingConversation", "GROUP", groupID).Return(false).Times(2)

	// Mock: Upsert sender's conversation
	mockConvRepo.On("Upsert", ctx, mock.MatchedBy(func(conv *models.Conversation) bool {
		return conv.UserID == senderID && conv.TargetID == groupID && conv.Type == "GROUP"
	})).Return(nil)

	// Mock: Increment unread for other members (now with 5 params including content)
	mockConvRepo.On("IncrementUnread", ctx, member1, "GROUP", groupID, content).Return(nil)
	mockConvRepo.On("IncrementUnread", ctx, member2, "GROUP", groupID, content).Return(nil)

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

	// B006: Also send to sender's other devices
	mockHub.On("SendToUser", senderID, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)
		return msg["type"] == "new_message"
	})).Return()

	// Execute
	msg, err := svc.SendGroupMessage(ctx, senderID, groupID, content)

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
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	nonMemberID := uuid.New()
	groupID := uuid.New()
	content := "I shouldn't be able to send this"

	// Mock: Sender is NOT a member
	mockGroupRepo.On("IsMember", ctx, groupID, nonMemberID).Return(false, nil)

	// Execute
	msg, err := svc.SendGroupMessage(ctx, nonMemberID, groupID, content)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, msg)
	assert.Equal(t, "sender is not a member of the group", err.Error())

	mockGroupRepo.AssertExpectations(t)
	// Message should NOT be created
	mockMsgRepo.AssertNotCalled(t, "Create")
}

func TestSendGroupMessage_BroadcastsToAllMembers(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

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
	mockGroupRepo.On("IsMember", ctx, groupID, senderID).Return(true, nil)
	mockMsgRepo.On("Create", ctx, mock.Anything).Return(nil)
	// Mock: Get sender info for response
	mockUserRepo.On("FindByID", ctx, senderID).Return(&models.User{
		BaseModel: models.BaseModel{ID: senderID},
		Username:  "Sender",
	}, nil)
	// Mock: Create Receipts [F06]
	// Mock: Create Receipts [F06] - Batch
	mockReceiptRepo.On("CreateBatch", ctx, mock.MatchedBy(func(receipts []*models.MessageReceipt) bool {
		return len(receipts) == 4
	})).Return(nil).Once()
	mockGroupRepo.On("GetMembers", ctx, groupID).Return(members, nil)
	mockConvRepo.On("Upsert", ctx, mock.Anything).Return(nil)

	// Mock: Check if viewing conversation (returns false for all members in this test)
	mockHub.On("IsUserViewingConversation", "GROUP", groupID).Return(false).Times(4)

	// Mock for each other member
	for _, memberID := range otherMemberIDs {
		mockConvRepo.On("IncrementUnread", ctx, memberID, "GROUP", groupID, mock.AnythingOfType("string")).Return(nil)
		mockHub.On("SendToUser", memberID, mock.Anything).Return()
	}
	// B006: Also send to sender's other devices
	mockHub.On("SendToUser", senderID, mock.Anything).Return()

	// Execute
	msg, err := svc.SendGroupMessage(ctx, senderID, groupID, "Test broadcast")

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
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	senderID := uuid.New()
	member1 := uuid.New()
	member2 := uuid.New()
	groupID := uuid.New()

	members := []models.GroupMember{
		{GroupID: groupID, UserID: senderID, Role: "MEMBER"},
		{GroupID: groupID, UserID: member1, Role: "MEMBER"},
		{GroupID: groupID, UserID: member2, Role: "MEMBER"},
	}

	mockGroupRepo.On("IsMember", ctx, groupID, senderID).Return(true, nil)
	mockMsgRepo.On("Create", ctx, mock.Anything).Return(nil)
	// Mock: Get sender info for response
	mockUserRepo.On("FindByID", ctx, senderID).Return(&models.User{
		BaseModel: models.BaseModel{ID: senderID},
		Username:  "Sender",
	}, nil)
	// Mock: Create Receipts [F06]
	// Mock: Create Receipts [F06] - Batch
	mockReceiptRepo.On("CreateBatch", ctx, mock.MatchedBy(func(receipts []*models.MessageReceipt) bool {
		return len(receipts) == 2
	})).Return(nil).Once()
	mockGroupRepo.On("GetMembers", ctx, groupID).Return(members, nil)

	// Mock: Check if viewing conversation (returns false for all members)
	mockHub.On("IsUserViewingConversation", "GROUP", groupID).Return(false).Times(2)

	// Expect conversation upsert for sender (unread = 0)
	mockConvRepo.On("Upsert", ctx, mock.MatchedBy(func(conv *models.Conversation) bool {
		return conv.UserID == senderID && conv.UnreadCount == 0
	})).Return(nil)

	// Expect increment unread for other members (with 5th param)
	mockConvRepo.On("IncrementUnread", ctx, member1, "GROUP", groupID, mock.AnythingOfType("string")).Return(nil)
	mockConvRepo.On("IncrementUnread", ctx, member2, "GROUP", groupID, mock.AnythingOfType("string")).Return(nil)

	mockHub.On("SendToUser", mock.Anything, mock.Anything).Return()
	mockHub.On("SendToUser", senderID, mock.Anything).Return()

	// Execute
	_, err := svc.SendGroupMessage(ctx, senderID, groupID, "Update conversations")

	// Assert
	assert.NoError(t, err)

	// Verify conversations were updated
	mockConvRepo.AssertCalled(t, "Upsert", ctx, mock.MatchedBy(func(conv *models.Conversation) bool {
		return conv.UserID == senderID
	}))
	mockConvRepo.AssertCalled(t, "IncrementUnread", ctx, member1, "GROUP", groupID, "Update conversations")
	mockConvRepo.AssertCalled(t, "IncrementUnread", ctx, member2, "GROUP", groupID, "Update conversations")
}

func TestSendGroupMessage_SenderDoesNotReceiveOwnMessage(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)

	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	senderID := uuid.New()
	member1 := uuid.New()
	groupID := uuid.New()

	members := []models.GroupMember{
		{GroupID: groupID, UserID: senderID, Role: "ADMIN"},
		{GroupID: groupID, UserID: member1, Role: "MEMBER"},
	}

	mockGroupRepo.On("IsMember", ctx, groupID, senderID).Return(true, nil)
	mockMsgRepo.On("Create", ctx, mock.Anything).Return(nil)
	// Mock: Get sender info for response
	mockUserRepo.On("FindByID", ctx, senderID).Return(&models.User{
		BaseModel: models.BaseModel{ID: senderID},
		Username:  "Sender",
	}, nil)
	// Mock: Create Receipts [F06]
	// Mock: Create Receipts [F06] - Batch
	mockReceiptRepo.On("CreateBatch", ctx, mock.MatchedBy(func(receipts []*models.MessageReceipt) bool {
		return len(receipts) == 1
	})).Return(nil).Once()
	mockGroupRepo.On("GetMembers", ctx, groupID).Return(members, nil)

	// Mock: Check if viewing conversation (returns false for all members)
	mockHub.On("IsUserViewingConversation", "GROUP", groupID).Return(false).Times(1)

	mockConvRepo.On("Upsert", ctx, mock.Anything).Return(nil)
	mockConvRepo.On("IncrementUnread", ctx, member1, "GROUP", groupID, mock.AnythingOfType("string")).Return(nil)
	mockHub.On("SendToUser", member1, mock.Anything).Return()

	// B006: Also send to sender's other devices
	mockHub.On("SendToUser", senderID, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)
		return msg["type"] == "new_message"
	})).Return()

	// Execute
	_, err := svc.SendGroupMessage(ctx, senderID, groupID, "Own message test")

	// Assert
	assert.NoError(t, err)

	// B006: Sender SHOULD receive the message (for multi-device sync)
	mockHub.AssertCalled(t, "SendToUser", senderID, mock.Anything)
	// And member1 should receive it
	mockHub.AssertCalled(t, "SendToUser", member1, mock.Anything)
}

func TestGetMessageReceipts_Success_Sender(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)
	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	userID := uuid.New()
	msgID := uuid.New()

	msg := &models.Message{
		BaseModel: models.BaseModel{ID: msgID},
		SenderID:  userID,
	}

	mockMsgRepo.On("FindByID", ctx, msgID).Return(msg, nil)

	mockReceipts := []models.MessageReceipt{{Status: "READ"}}
	mockReceiptRepo.On("FindByMessageID", ctx, msgID).Return(mockReceipts, nil)

	res, err := svc.GetMessageReceipts(ctx, userID, msgID)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestGetMessageReceipts_Forbidden(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockHub := new(MockHub)
	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	userID := uuid.New()
	otherUser := uuid.New()
	msgID := uuid.New()

	msg := &models.Message{
		BaseModel: models.BaseModel{ID: msgID},
		SenderID:  otherUser, // Not sender
		// No ReceiverID, No GroupID -> generic forbidden
	}

	mockMsgRepo.On("FindByID", ctx, msgID).Return(msg, nil)

	_, err := svc.GetMessageReceipts(ctx, userID, msgID)
	assert.Error(t, err)
	assert.Equal(t, "access denied", err.Error())
}
