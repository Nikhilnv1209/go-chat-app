package handlers_test

import (
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

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepository) UpdateOnlineStatus(userID uuid.UUID, isOnline bool, lastSeen time.Time) error {
	args := m.Called(userID, isOnline, lastSeen)
	return args.Error(0)
}

func setupWSTest() (*handlers.WSHandler, *MockAuthService, *MockUserRepository, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockAuthService := new(MockAuthService)
	mockUserRepo := new(MockUserRepository)

	hub := websocket.NewHub(mockUserRepo)
	go hub.Run() // Start hub

	handler := handlers.NewWSHandler(hub, mockAuthService)
	r := gin.New()
	return handler, mockAuthService, mockUserRepo, r
}

func TestServeWS_NoToken(t *testing.T) {
	handler, _, _, r := setupWSTest()
	r.GET("/ws", handler.ServeWS)

	req, _ := http.NewRequest("GET", "/ws", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestServeWS_InvalidToken(t *testing.T) {
	handler, mockService, _, r := setupWSTest()
	r.GET("/ws", handler.ServeWS)

	mockService.On("ValidateToken", "invalid_token").Return(uuid.Nil, errors.ErrUnauthorized)

	req, _ := http.NewRequest("GET", "/ws?token=invalid_token", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestServeWS_Success(t *testing.T) {
	handler, mockService, mockRepo, r := setupWSTest()
	r.GET("/ws", handler.ServeWS)

	userID := uuid.New()
	mockService.On("ValidateToken", "valid_token").Return(userID, nil)

	// Expect UpdateOnlineStatus to be called with true (Online)
	mockRepo.On("UpdateOnlineStatus", userID, true, mock.Anything).Return(nil).Maybe()

	// Expect UpdateOnlineStatus to be called with false (Offline) eventually
	mockRepo.On("UpdateOnlineStatus", userID, false, mock.Anything).Return(nil).Maybe()

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
