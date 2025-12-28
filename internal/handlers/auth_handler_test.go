package handlers_test

import (
	"bytes"
	"chat-app/internal/errors"
	"chat-app/internal/handlers"
	"chat-app/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(username, email, password string) (string, string, *models.User, error) {
	args := m.Called(username, email, password)
	if args.Get(2) == nil {
		return args.String(0), args.String(1), nil, args.Error(3)
	}
	return args.String(0), args.String(1), args.Get(2).(*models.User), args.Error(3)
}

func (m *MockAuthService) Login(email, password string) (string, string, *models.User, error) {
	args := m.Called(email, password)
	if args.Get(2) == nil {
		return args.String(0), args.String(1), nil, args.Error(3)
	}
	return args.String(0), args.String(1), args.Get(2).(*models.User), args.Error(3)
}

func (m *MockAuthService) Refresh(refreshToken string) (string, string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) Logout(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	args := m.Called(tokenString)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAuthService) SearchUsers(query string, excludeUserID uuid.UUID) ([]models.User, error) {
	args := m.Called(query, excludeUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockAuthService) GetUser(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func setupAuthTest() (*handlers.AuthHandler, *MockAuthService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	handler := handlers.NewAuthHandler(mockService)
	r := gin.New()
	return handler, mockService, r
}

func TestRegister_Success(t *testing.T) {
	handler, mockService, r := setupAuthTest()
	r.POST("/register", handler.Register)

	user := &models.User{Username: "test", Email: "test@example.com"}
	mockService.On("Register", "test", "test@example.com", "password123").Return("test_token", "refresh_token", user, nil)

	body := `{"username":"test", "email":"test@example.com", "password":"password123"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	handler, mockService, r := setupAuthTest()
	r.POST("/register", handler.Register)

	mockService.On("Register", "test", "existing@example.com", "password123").Return("", "", nil, errors.ErrEmailExists)

	body := `{"username":"test", "email":"existing@example.com", "password":"password123"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegister_InvalidInput(t *testing.T) {
	handler, _, r := setupAuthTest()
	r.POST("/register", handler.Register)

	// Missing password
	body := `{"username":"test", "email":"test@example.com"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	handler, mockService, r := setupAuthTest()
	r.POST("/login", handler.Login)

	user := &models.User{Email: "test@example.com"}
	mockService.On("Login", "test@example.com", "password123").Return("valid_token", "refresh_token", user, nil)

	body := `{"email":"test@example.com", "password":"password123"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "valid_token", response["token"])
}

func TestLogin_InvalidCredentials(t *testing.T) {
	handler, mockService, r := setupAuthTest()
	r.POST("/login", handler.Login)

	mockService.On("Login", "test@example.com", "wrongpass").Return("", "", nil, errors.ErrInvalidCredentials)

	body := `{"email":"test@example.com", "password":"wrongpass"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
