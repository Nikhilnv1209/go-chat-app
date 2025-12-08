package handlers_test

import (
	"bytes"
	"chat-app/internal/handlers"
	"chat-app/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGroupService
type MockGroupService struct {
	mock.Mock
}

func (m *MockGroupService) Create(creatorID uuid.UUID, name string, memberIDs []uuid.UUID) (*models.Group, error) {
	args := m.Called(creatorID, name, memberIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *MockGroupService) AddMember(adminID, groupID, newMemberID uuid.UUID) error {
	args := m.Called(adminID, groupID, newMemberID)
	return args.Error(0)
}

func (m *MockGroupService) RemoveMember(adminID, groupID, memberID uuid.UUID) error {
	args := m.Called(adminID, groupID, memberID)
	return args.Error(0)
}

func setupGroupTest() (*handlers.GroupHandler, *MockGroupService, *MockAuthService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockGroupService := new(MockGroupService)
	mockAuthService := new(MockAuthService)
	handler := handlers.NewGroupHandler(mockGroupService, mockAuthService)
	r := gin.New()
	return handler, mockGroupService, mockAuthService, r
}

func TestCreateGroup_Success(t *testing.T) {
	handler, mockGroupService, mockAuthService, r := setupGroupTest()
	r.POST("/groups", handler.CreateGroup)

	creatorID := uuid.New()
	member1 := uuid.New()
	member2 := uuid.New()
	groupID := uuid.New()

	// Mock token validation
	mockAuthService.On("ValidateToken", "valid_token").Return(creatorID, nil)

	// Mock group creation
	group := &models.Group{
		BaseModel: models.BaseModel{ID: groupID},
		Name:      "Test Group",
	}
	mockGroupService.On("Create", creatorID, "Test Group", []uuid.UUID{member1, member2}).Return(group, nil)

	// Prepare request
	body := map[string]interface{}{
		"name":       "Test Group",
		"member_ids": []string{member1.String(), member2.String()},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer valid_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, groupID.String(), response["id"])
	assert.Equal(t, "Test Group", response["name"])

	mockAuthService.AssertExpectations(t)
	mockGroupService.AssertExpectations(t)
}

func TestCreateGroup_Unauthorized_NoToken(t *testing.T) {
	handler, _, _, r := setupGroupTest()
	r.POST("/groups", handler.CreateGroup)

	body := `{"name": "Test Group", "member_ids": []}`
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateGroup_Unauthorized_InvalidToken(t *testing.T) {
	handler, _, mockAuthService, r := setupGroupTest()
	r.POST("/groups", handler.CreateGroup)

	// Mock invalid token
	mockAuthService.On("ValidateToken", "invalid_token").Return(uuid.Nil, assert.AnError)

	body := `{"name": "Test Group", "member_ids": []}`
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer invalid_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateGroup_BadRequest_MissingName(t *testing.T) {
	handler, _, mockAuthService, r := setupGroupTest()
	r.POST("/groups", handler.CreateGroup)

	creatorID := uuid.New()
	mockAuthService.On("ValidateToken", "valid_token").Return(creatorID, nil)

	// Missing name field
	body := `{"member_ids": []}`
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer valid_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddMember_Success(t *testing.T) {
	handler, mockGroupService, mockAuthService, r := setupGroupTest()
	r.POST("/groups/:id/members", handler.AddMember)

	adminID := uuid.New()
	groupID := uuid.New()
	newMemberID := uuid.New()

	// Mock token validation
	mockAuthService.On("ValidateToken", "admin_token").Return(adminID, nil)

	// Mock adding member
	mockGroupService.On("AddMember", adminID, groupID, newMemberID).Return(nil)

	// Prepare request
	body := map[string]string{"user_id": newMemberID.String()}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/groups/"+groupID.String()+"/members", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer admin_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "member added successfully", response["message"])

	mockAuthService.AssertExpectations(t)
	mockGroupService.AssertExpectations(t)
}

func TestAddMember_Forbidden_NotAdmin(t *testing.T) {
	handler, mockGroupService, mockAuthService, r := setupGroupTest()
	r.POST("/groups/:id/members", handler.AddMember)

	regularUserID := uuid.New()
	groupID := uuid.New()
	newMemberID := uuid.New()

	// Mock token validation
	mockAuthService.On("ValidateToken", "user_token").Return(regularUserID, nil)

	// Mock service returning forbidden error
	forbiddenErr := errors.New("only admins can add members")
	mockGroupService.On("AddMember", regularUserID, groupID, newMemberID).Return(forbiddenErr)

	// Prepare request
	body := map[string]string{"user_id": newMemberID.String()}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/groups/"+groupID.String()+"/members", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer user_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert - Should be 403 Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)

	mockAuthService.AssertExpectations(t)
	mockGroupService.AssertExpectations(t)
}

func TestAddMember_BadRequest_InvalidGroupID(t *testing.T) {
	handler, _, mockAuthService, r := setupGroupTest()
	r.POST("/groups/:id/members", handler.AddMember)

	adminID := uuid.New()
	mockAuthService.On("ValidateToken", "admin_token").Return(adminID, nil)

	body := `{"user_id": "` + uuid.New().String() + `"}`
	req, _ := http.NewRequest("POST", "/groups/invalid-uuid/members", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer admin_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddMember_BadRequest_MissingUserID(t *testing.T) {
	handler, _, mockAuthService, r := setupGroupTest()
	r.POST("/groups/:id/members", handler.AddMember)

	adminID := uuid.New()
	groupID := uuid.New()

	mockAuthService.On("ValidateToken", "admin_token").Return(adminID, nil)

	// Missing user_id
	body := `{}`
	req, _ := http.NewRequest("POST", "/groups/"+groupID.String()+"/members", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer admin_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
