package handlers

import (
	"context"
	"encoding/json"
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

// MockMessageService
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) SendDirectMessage(ctx context.Context, senderID, receiverID uuid.UUID, content string) (*models.Message, error) {
	args := m.Called(ctx, senderID, receiverID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockMessageService) SendGroupMessage(ctx context.Context, senderID, groupID uuid.UUID, content string) (*models.Message, error) {
	args := m.Called(ctx, senderID, groupID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockMessageService) GetHistory(ctx context.Context, userID, targetID uuid.UUID, convType string, limit int, beforeID *uuid.UUID) ([]models.Message, error) {
	args := m.Called(ctx, userID, targetID, convType, limit, beforeID)
	return args.Get(0).([]models.Message), args.Error(1)
}

func (m *MockMessageService) MarkAsRead(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) error {
	args := m.Called(ctx, userID, messageIDs)
	return args.Error(0)
}

func (m *MockMessageService) MarkAsDelivered(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) error {
	args := m.Called(ctx, userID, messageIDs)
	return args.Error(0)
}

func (m *MockMessageService) GetMessageReceipts(ctx context.Context, userID, messageID uuid.UUID) ([]models.MessageReceipt, error) {
	args := m.Called(ctx, userID, messageID)
	return args.Get(0).([]models.MessageReceipt), args.Error(1)
}

func (m *MockMessageService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockMessageService) BroadcastTypingIndicator(ctx context.Context, userID uuid.UUID, username, convType string, targetID uuid.UUID, isTyping bool) error {
	args := m.Called(ctx, userID, username, convType, targetID, isTyping)
	return args.Error(0)
}

// Helper middleware to set userID in context (simulates what AuthMiddleware does)
func mockAuthMiddleware(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

// Tests

func TestGetConversations_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

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
	mockConvRepo.On("FindByUser", mock.AnythingOfType("*context.timerCtx"), userID).Return(conversations, nil)
	mockUserRepo.On("FindByID", mock.AnythingOfType("*context.timerCtx"), targetUserID).Return(&models.User{
		BaseModel: models.BaseModel{ID: targetUserID},
		Username:  "Bob",
	}, nil)
	mockGroupRepo.On("FindByID", mock.AnythingOfType("*context.timerCtx"), groupID).Return(&models.Group{
		BaseModel: models.BaseModel{ID: groupID},
		Name:      "Family",
	}, nil)
	// Mock GetMembers for group member count
	mockGroupRepo.On("GetMembers", mock.AnythingOfType("*context.timerCtx"), groupID).Return([]models.GroupMember{
		{GroupID: groupID, UserID: userID, Role: "MEMBER"},
		{GroupID: groupID, UserID: uuid.New(), Role: "MEMBER"},
	}, nil)

	// Create router with mock auth middleware
	r := gin.New()
	r.GET("/conversations", mockAuthMiddleware(userID), handler.GetConversations)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/conversations", nil)
	r.ServeHTTP(w, req)

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
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

	// Create router WITHOUT auth middleware - userID won't be in context
	r := gin.New()
	r.GET("/conversations", handler.GetConversations)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/conversations", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetMessages_DM_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

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
	mockMsgService.On("GetHistory", mock.AnythingOfType("*context.timerCtx"), userID, targetID, mock.Anything, 50, (*uuid.UUID)(nil)).Return(messages, nil)
	mockConvRepo.On("ResetUnread", mock.AnythingOfType("*context.timerCtx"), userID, "DM", targetID).Return(nil)

	// Create router with mock auth middleware
	r := gin.New()
	r.GET("/messages", mockAuthMiddleware(userID), handler.GetMessages)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages?target_id="+targetID.String()+"&type=DM", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)

	mockMsgService.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
}

func TestGetMessages_Group_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

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
	mockGroupRepo.On("IsMember", mock.AnythingOfType("*context.timerCtx"), groupID, userID).Return(true, nil)
	mockMsgService.On("GetHistory", mock.AnythingOfType("*context.timerCtx"), userID, groupID, mock.Anything, 50, (*uuid.UUID)(nil)).Return(messages, nil)
	mockConvRepo.On("ResetUnread", mock.AnythingOfType("*context.timerCtx"), userID, "GROUP", groupID).Return(nil)

	// Create router with mock auth middleware
	r := gin.New()
	r.GET("/messages", mockAuthMiddleware(userID), handler.GetMessages)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages?target_id="+groupID.String()+"&type=GROUP", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)

	mockGroupRepo.AssertExpectations(t)
	mockMsgService.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
}

func TestGetMessages_Group_NotMember(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

	// Test data
	userID := uuid.New()
	groupID := uuid.New()

	// Mock expectations
	mockGroupRepo.On("IsMember", mock.AnythingOfType("*context.timerCtx"), groupID, userID).Return(false, nil)

	// Create router with mock auth middleware
	r := gin.New()
	r.GET("/messages", mockAuthMiddleware(userID), handler.GetMessages)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages?target_id="+groupID.String()+"&type=GROUP", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)

	mockGroupRepo.AssertExpectations(t)
}

func TestGetMessages_MissingTargetID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

	// Test data
	userID := uuid.New()

	// Create router with mock auth middleware
	r := gin.New()
	r.GET("/messages", mockAuthMiddleware(userID), handler.GetMessages)

	// Create request without target_id
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages?type=DM", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMessages_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

	// Test data
	userID := uuid.New()
	targetID := uuid.New()

	// Create router with mock auth middleware
	r := gin.New()
	r.GET("/messages", mockAuthMiddleware(userID), handler.GetMessages)

	// Create request with invalid type
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages?target_id="+targetID.String()+"&type=INVALID", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMessages_CustomLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

	// Test data
	userID := uuid.New()
	targetID := uuid.New()

	messages := []models.Message{}

	// Mock expectations - should use limit=20
	mockMsgService.On("GetHistory", mock.AnythingOfType("*context.timerCtx"), userID, targetID, mock.Anything, 20, (*uuid.UUID)(nil)).Return(messages, nil)
	mockConvRepo.On("ResetUnread", mock.AnythingOfType("*context.timerCtx"), userID, "DM", targetID).Return(nil)

	// Create router with mock auth middleware
	r := gin.New()
	r.GET("/messages", mockAuthMiddleware(userID), handler.GetMessages)

	// Create request with custom limit
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages?target_id="+targetID.String()+"&type=DM&limit=20", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	mockMsgService.AssertExpectations(t)
	mockConvRepo.AssertExpectations(t)
}

func TestMarkRead_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

	// Test data
	userID := uuid.New()
	msgID := uuid.New()

	// Mock expectations
	mockMsgService.On("MarkAsRead", mock.AnythingOfType("*context.timerCtx"), userID, []uuid.UUID{msgID}).Return(nil)

	// Create router with mock auth middleware
	r := gin.New()
	r.POST("/messages/:id/read", mockAuthMiddleware(userID), handler.MarkRead)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages/"+msgID.String()+"/read", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	mockMsgService.AssertExpectations(t)
}

func TestMarkRead_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

	// Test data
	userID := uuid.New()

	// Create router with mock auth middleware
	r := gin.New()
	r.POST("/messages/:id/read", mockAuthMiddleware(userID), handler.MarkRead)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages/invalid-uuid/read", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetReceipts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

	userID := uuid.New()
	messageID := uuid.New()

	mockReceipts := []models.MessageReceipt{
		{
			BaseModel: models.BaseModel{ID: uuid.New(), CreatedAt: time.Now()},
			MessageID: messageID,
			UserID:    uuid.New(),
			Status:    "READ",
		},
	}
	mockMsgService.On("GetMessageReceipts", mock.AnythingOfType("*context.timerCtx"), userID, messageID).Return(mockReceipts, nil)

	// Create router with mock auth middleware
	r := gin.New()
	r.GET("/messages/:id/receipts", mockAuthMiddleware(userID), handler.GetReceipts)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages/"+messageID.String()+"/receipts", nil)
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
	gin.SetMode(gin.TestMode)

	mockConvRepo := new(MockConversationRepo)
	mockMsgRepo := new(MockMessageRepo)
	mockUserRepo := new(MockUserRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockMsgService := new(MockMessageService)

	handler := NewChatHandler(mockConvRepo, mockMsgRepo, mockUserRepo, mockGroupRepo, mockMsgService)

	userID := uuid.New()

	// Create router with mock auth middleware
	r := gin.New()
	r.GET("/messages/:id/receipts", mockAuthMiddleware(userID), handler.GetReceipts)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages/invalid-uuid/receipts", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
