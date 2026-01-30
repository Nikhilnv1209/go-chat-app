package handlers_test

import (
	"context"
	"chat-app/internal/errors"
	"chat-app/internal/handlers"
	"chat-app/internal/models"
	"chat-app/internal/websocket"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository to mock database updates in Hub
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepository) UpdateOnlineStatus(ctx context.Context, userID uuid.UUID, isOnline bool, lastSeen time.Time) error {
	args := m.Called(ctx, userID, isOnline, lastSeen)
	return args.Error(0)
}

func (m *MockUserRepository) Search(ctx context.Context, query string, excludeUserID uuid.UUID) ([]models.User, error) {
	args := m.Called(ctx, query, excludeUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

// MockConversationRepository to mock conversation lookups for presence broadcasting
type MockConversationRepository struct {
	mock.Mock
}

func (m *MockConversationRepository) Upsert(ctx context.Context, conv *models.Conversation) error {
	args := m.Called(ctx, conv)
	return args.Error(0)
}

func (m *MockConversationRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Conversation), args.Error(1)
}

func (m *MockConversationRepository) IncrementUnread(ctx context.Context, userID uuid.UUID, convType string, targetID uuid.UUID, lastMessage string) error {
	args := m.Called(ctx, userID, convType, targetID, lastMessage)
	return args.Error(0)
}

func (m *MockConversationRepository) ResetUnread(ctx context.Context, userID uuid.UUID, convType string, targetID uuid.UUID) error {
	args := m.Called(ctx, userID, convType, targetID)
	return args.Error(0)
}

func (m *MockConversationRepository) FindContactsOfUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

func setupWSTest() (*handlers.WSHandler, *MockAuthService, *MockUserRepository, *MockConversationRepository, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockAuthService := new(MockAuthService)
	mockUserRepo := new(MockUserRepository)
	mockConvRepo := new(MockConversationRepository)

	hub := websocket.NewHub(mockUserRepo, mockConvRepo)
	go hub.Run() // Start hub

	handler := handlers.NewWSHandler(hub, mockAuthService)
	r := gin.New()
	return handler, mockAuthService, mockUserRepo, mockConvRepo, r
}

func TestServeWS_NoToken(t *testing.T) {
	handler, _, _, _, r := setupWSTest()
	r.GET("/ws", handler.ServeWS)

	req, _ := http.NewRequest("GET", "/ws", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestServeWS_InvalidToken(t *testing.T) {
	handler, mockService, _, _, r := setupWSTest()
	r.GET("/ws", handler.ServeWS)

	mockService.On("ValidateToken", "invalid_token").Return(uuid.Nil, errors.ErrUnauthorized)

	req, _ := http.NewRequest("GET", "/ws?token=invalid_token", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestServeWS_Success(t *testing.T) {
	handler, mockService, mockRepo, mockConvRepo, r := setupWSTest()
	r.GET("/ws", handler.ServeWS)

	userID := uuid.New()
	mockService.On("ValidateToken", "valid_token").Return(userID, nil)

	// Expect UpdateOnlineStatus to be called with true (Online)
	mockRepo.On("UpdateOnlineStatus", mock.AnythingOfType("*context.timerCtx"), userID, true, mock.Anything).Return(nil).Maybe()

	// Expect UpdateOnlineStatus to be called with false (Offline) eventually
	mockRepo.On("UpdateOnlineStatus", mock.AnythingOfType("*context.timerCtx"), userID, false, mock.Anything).Return(nil).Maybe()

	// Expect FindByID to be called for getting user info (for presence broadcasting)
	mockRepo.On("FindByID", mock.AnythingOfType("*context.timerCtx"), userID).Return(&models.User{
		BaseModel: models.BaseModel{ID: userID},
		Username:  "TestUser",
		IsOnline:  true,
	}, nil).Maybe()

	// Expect FindByUser to be called for initial presence sync
	mockConvRepo.On("FindByUser", mock.AnythingOfType("*context.timerCtx"), userID).Return([]models.Conversation{}, nil).Maybe()

	// Expect FindContactsOfUser to be called for presence broadcasting
	mockConvRepo.On("FindContactsOfUser", mock.AnythingOfType("*context.timerCtx"), userID).Return([]uuid.UUID{}, nil).Maybe()

	// Create a test server to handle the websocket upgrade
	s := httptest.NewServer(r)
	defer s.Close()

	// Convert http URL to ws URL
	u := "ws" + s.URL[4:] + "/ws?token=valid_token"

	// Connect to the server
	ws, _, err := gorilla.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)

	// Optionally send a ping or message to verify
	err = ws.WriteMessage(gorilla.PingMessage, []byte{})
	assert.NoError(t, err)

	// Give it a moment to process the connection (and async DB call)
	time.Sleep(50 * time.Millisecond)

	ws.Close()
	// Give it a moment to process disconnection
	time.Sleep(50 * time.Millisecond)
}
