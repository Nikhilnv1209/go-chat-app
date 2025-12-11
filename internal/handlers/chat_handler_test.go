package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"chat-app/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories and services
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

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) UpdateOnlineStatus(userID uuid.UUID, isOnline bool, lastSeen time.Time) error {
	args := m.Called(userID, isOnline, lastSeen)
	return args.Error(0)
}

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

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(username, email, password string) (*models.User, error) {
	args := m.Called(username, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) Login(email, password string) (string, *models.User, error) {
	args := m.Called(email, password)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*models.User), args.Error(2)
}

func (m *MockAuthService) ValidateToken(token string) (uuid.UUID, error) {
	args := m.Called(token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

// MockMessageService
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) SendDirectMessage(senderID, receiverID uuid.UUID, content string) (*models.Message, error) {
	args := m.Called(senderID, receiverID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockMessageService) SendGroupMessage(senderID, groupID uuid.UUID, content string) (*models.Message, error) {
	args := m.Called(senderID, groupID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockMessageService) GetHistory(userID, targetID uuid.UUID, convType string, limit, beforeID int) ([]models.Message, error) {
	args := m.Called(userID, targetID, convType, limit, beforeID)
	return args.Get(0).([]models.Message), args.Error(1)
}

func (m *MockMessageService) MarkAsRead(userID uuid.UUID, messageIDs []uuid.UUID) error {
	args := m.Called(userID, messageIDs)
	return args.Error(0)
}

func (m *MockMessageService) MarkAsDelivered(userID uuid.UUID, messageIDs []uuid.UUID) error {
	args := m.Called(userID, messageIDs)
	return args.Error(0)
}

func (m *MockMessageService) GetMessageReceipts(userID, messageID uuid.UUID) ([]models.MessageReceipt, error) {
	args := m.Called(userID, messageID)
	return args.Get(0).([]models.MessageReceipt), args.Error(1)
}

func (m *MockMessageService) GetUserInfo(userID uuid.UUID) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockMessageService) BroadcastTypingIndicator(userID uuid.UUID, username, convType string, targetID uuid.UUID, isTyping bool) error {
	args := m.Called(userID, username, convType, targetID, isTyping)
	return args.Error(0)
}

// Tests

func TestGetConversations_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Test data
	userID := uuid.New()
	targetUserID := uuid.New()
	groupID := uuid.New()
	now := time.Now()

	conversations := []models.Conversation{
		{
			BaseModel:     models.BaseModel{ID: uuid.New()},
			UserID:        userID,
			Type:          "DM",
			TargetID:      targetUserID,
			LastMessageAt: now,
			UnreadCount:   3,
		},
		{
			BaseModel:     models.BaseModel{ID: uuid.New()},
			UserID:        userID,
			Type:          "GROUP",
			TargetID:      groupID,
			LastMessageAt: now.Add(-1 * time.Hour),
			UnreadCount:   0,
		},
	}

	// Mock expectations
	mockAuthService.On("ValidateToken", "test-token").Return(userID, nil)
	mockConvRepo.On("FindByUser", userID).Return(conversations, nil)
	mockUserRepo.On("FindByID", targetUserID).Return(&models.User{
		BaseModel: models.BaseModel{ID: targetUserID},
		Username:  "Bob",
	}, nil)
	mockGroupRepo.On("FindByID", groupID).Return(&models.Group{
		BaseModel: models.BaseModel{ID: groupID},
		Name:      "Family",
	}, nil)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/conversations", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")

	// Execute
	handler.GetConversations(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []ConversationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "DM", response[0].Type)
	assert.Equal(t, "Bob", response[0].TargetName)
	assert.Equal(t, 3, response[0].UnreadCount)
	assert.Equal(t, "GROUP", response[1].Type)
	assert.Equal(t, "Family", response[1].TargetName)

	mockAuthService.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockGroupRepo.AssertExpectations(t)
}

func TestGetConversations_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Create request without auth header
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/conversations", nil)

	// Execute
	handler.GetConversations(c)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetConversations_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Mock expectations
	mockAuthService.On("ValidateToken", "invalid-token").Return(uuid.Nil, errors.New("invalid token"))

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/conversations", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	// Execute
	handler.GetConversations(c)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuthService.AssertExpectations(t)
}

func TestGetMessages_DM_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Test data
	userID := uuid.New()
	targetID := uuid.New()

	messages := []models.Message{
		{
			BaseModel:  models.BaseModel{ID: uuid.New()},
			SenderID:   targetID,
			ReceiverID: &userID,
			Content:    "Hello",
		},
		{
			BaseModel:  models.BaseModel{ID: uuid.New()},
			SenderID:   userID,
			ReceiverID: &targetID,
			Content:    "Hey!",
		},
	}

	// Mock expectations
	mockAuthService.On("ValidateToken", "test-token").Return(userID, nil)
	// Change to use Service mock instead of Repo mock
	mockMsgService.On("GetHistory", userID, targetID, "DM", 50, 0).Return(messages, nil)
	// mockConvRepo.On("ResetUnread"...) removed because it's assumed handled by Service or separate logic if not in Service
	// But in ChatHandler.GetMessages, we still call convRepo.ResetUnread manually unless we moved that too.
	// Oh right, I kept ResetUnread inside ChatHandler in my "replace_file_content" logic earlier?
	// Let's re-read the code logic.
	// I replaced h.msgRepo.FindByConversation with h.msgService.GetHistory.
	// But I LEFT h.convRepo.ResetUnread(userID, msgType, targetID).
	// So I still need to mock ResetUnread.
	mockConvRepo.On("ResetUnread", userID, "DM", targetID).Return(nil)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/messages?target_id="+targetID.String()+"&type=DM", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")

	// Execute
	handler.GetMessages(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)

	mockAuthService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
	mockMsgService.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
}

func TestGetMessages_Group_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Test data
	userID := uuid.New()
	groupID := uuid.New()

	messages := []models.Message{
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			SenderID:  userID,
			GroupID:   &groupID,
			Content:   "Hello group!",
		},
	}

	// Mock expectations
	// Mock expectations
	mockAuthService.On("ValidateToken", "test-token").Return(userID, nil)
	mockGroupRepo.On("IsMember", groupID, userID).Return(true, nil)
	mockMsgService.On("GetHistory", userID, groupID, "GROUP", 50, 0).Return(messages, nil)
	mockConvRepo.On("ResetUnread", userID, "GROUP", groupID).Return(nil)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/messages?target_id="+groupID.String()+"&type=GROUP", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")

	// Execute
	handler.GetMessages(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)

	mockAuthService.AssertExpectations(t)
	mockGroupRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
	mockMsgService.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
}

func TestGetMessages_Group_NotMember(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Test data
	userID := uuid.New()
	groupID := uuid.New()

	// Mock expectations
	mockAuthService.On("ValidateToken", "test-token").Return(userID, nil)
	mockGroupRepo.On("IsMember", groupID, userID).Return(false, nil)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/messages?target_id="+groupID.String()+"&type=GROUP", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")

	// Execute
	handler.GetMessages(c)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)

	mockAuthService.AssertExpectations(t)
	mockGroupRepo.AssertExpectations(t)
}

func TestGetMessages_MissingTargetID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Test data
	userID := uuid.New()

	// Mock expectations
	mockAuthService.On("ValidateToken", "test-token").Return(userID, nil)

	// Create request without target_id
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/messages?type=DM", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")

	// Execute
	handler.GetMessages(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockAuthService.AssertExpectations(t)
}

func TestGetMessages_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Test data
	userID := uuid.New()
	targetID := uuid.New()

	// Mock expectations
	mockAuthService.On("ValidateToken", "test-token").Return(userID, nil)

	// Create request with invalid type
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/messages?target_id="+targetID.String()+"&type=INVALID", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")

	// Execute
	handler.GetMessages(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockAuthService.AssertExpectations(t)
}

func TestGetMessages_CustomLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Test data
	userID := uuid.New()
	targetID := uuid.New()

	messages := []models.Message{}

	// Mock expectations - should use limit=20
	mockAuthService.On("ValidateToken", "test-token").Return(userID, nil)
	mockMsgService.On("GetHistory", userID, targetID, "DM", 20, 0).Return(messages, nil)
	mockConvRepo.On("ResetUnread", userID, "DM", targetID).Return(nil)

	// Create request with custom limit
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/messages?target_id="+targetID.String()+"&type=DM&limit=20", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")

	// Execute
	handler.GetMessages(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	mockAuthService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
	mockMsgService.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
}

func TestMarkRead_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Test data
	userID := uuid.New()
	msgID := uuid.New()

	// Mock expectations
	mockAuthService.On("ValidateToken", "test-token").Return(userID, nil)
	mockMsgService.On("MarkAsRead", userID, []uuid.UUID{msgID}).Return(nil)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/messages/"+msgID.String()+"/read", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")
	// Set Gin param manually since router isn't matching path parameters in standalone context
	c.Params = gin.Params{{Key: "id", Value: msgID.String()}}

	// Execute
	handler.MarkRead(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	mockAuthService.AssertExpectations(t)
	mockMsgService.AssertExpectations(t)
}

func TestMarkRead_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)

	// Test data
	userID := uuid.New()

	// Mock expectations
	mockAuthService.On("ValidateToken", "test-token").Return(userID, nil)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/messages/invalid-uuid/read", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	// Execute
	handler.MarkRead(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockAuthService.AssertExpectations(t)
}

func TestGetReceipts_Success(t *testing.T) {
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)
	// Inject router
	r := gin.Default()
	r.GET("/messages/:id/receipts", handler.GetReceipts)

	userID := uuid.New()
	messageID := uuid.New()
	token := "valid_token"

	// Expectations
	mockAuthService.On("ValidateToken", token).Return(userID, nil)

	mockReceipts := []models.MessageReceipt{
		{
			BaseModel: models.BaseModel{ID: uuid.New(), CreatedAt: time.Now()},
			MessageID: messageID,
			UserID:    uuid.New(),
			Status:    "READ",
		},
	}
	mockMsgService.On("GetMessageReceipts", userID, messageID).Return(mockReceipts, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages/"+messageID.String()+"/receipts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []models.MessageReceipt
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "READ", resp[0].Status)
}

func TestGetReceipts_InvalidID(t *testing.T) {
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockAuthService := new(MockAuthService)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockAuthService, mockMsgService)
	r := gin.Default()
	r.GET("/messages/:id/receipts", handler.GetReceipts)

	mockAuthService.On("ValidateToken", "valid").Return(uuid.New(), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages/invalid-uuid/receipts", nil)
	req.Header.Set("Authorization", "Bearer valid")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
