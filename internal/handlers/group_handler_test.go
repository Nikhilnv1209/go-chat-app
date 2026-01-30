package handlers_test

import (
	"bytes"
	"context"
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

func (m *MockGroupService) Create(ctx context.Context, creatorID uuid.UUID, name string, memberIDs []uuid.UUID) (*models.Group, error) {
	args := m.Called(ctx, creatorID, name, memberIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *MockGroupService) AddMember(ctx context.Context, adminID, groupID, newMemberID uuid.UUID) error {
	args := m.Called(ctx, adminID, groupID, newMemberID)
	return args.Error(0)
}

func (m *MockGroupService) RemoveMember(ctx context.Context, adminID, groupID, memberID uuid.UUID) error {
	args := m.Called(ctx, adminID, groupID, memberID)
	return args.Error(0)
}

func setupGroupTest() (*handlers.GroupHandler, *MockGroupService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockGroupService := new(MockGroupService)
	handler := handlers.NewGroupHandler(mockGroupService)
	r := gin.New()
	return handler, mockGroupService, r
}

// Helper middleware to set userID in context (simulates what AuthMiddleware does)
func mockAuthMiddleware(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func TestCreateGroup_Success(t *testing.T) {
	handler, mockGroupService, r := setupGroupTest()

	creatorID := uuid.New()
	member1 := uuid.New()
	member2 := uuid.New()
	groupID := uuid.New()

	// Apply mock auth middleware that sets userID
	r.POST("/groups", mockAuthMiddleware(creatorID), handler.CreateGroup)

	// Mock group creation
	group := &models.Group{
		BaseModel: models.BaseModel{ID: groupID},
		Name:      "Test Group",
	}
	mockGroupService.On("Create", mock.AnythingOfType("*context.timerCtx"), creatorID, "Test Group", []uuid.UUID{member1, member2}).Return(group, nil)

	// Prepare request
	body := map[string]interface{}{
		"name":       "Test Group",
		"member_ids": []string{member1.String(), member2.String()},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, groupID.String(), response["id"])
	assert.Equal(t, "Test Group", response["name"])

	mockGroupService.AssertExpectations(t)
}

func TestCreateGroup_Unauthorized_NoUserInContext(t *testing.T) {
	handler, _, r := setupGroupTest()
	// No auth middleware - userID won't be in context
	r.POST("/groups", handler.CreateGroup)

	body := `{"name": "Test Group", "member_ids": []}`
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateGroup_BadRequest_MissingName(t *testing.T) {
	handler, _, r := setupGroupTest()

	creatorID := uuid.New()
	r.POST("/groups", mockAuthMiddleware(creatorID), handler.CreateGroup)

	// Missing name field
	body := `{"member_ids": []}`
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddMember_Success(t *testing.T) {
	handler, mockGroupService, r := setupGroupTest()

	adminID := uuid.New()
	groupID := uuid.New()
	newMemberID := uuid.New()

	r.POST("/groups/:id/members", mockAuthMiddleware(adminID), handler.AddMember)

	// Mock adding member
	mockGroupService.On("AddMember", mock.AnythingOfType("*context.timerCtx"), adminID, groupID, newMemberID).Return(nil)

	// Prepare request
	body := map[string]string{"user_id": newMemberID.String()}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/groups/"+groupID.String()+"/members", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "member added successfully", response["message"])

	mockGroupService.AssertExpectations(t)
}

func TestAddMember_Forbidden_NotAdmin(t *testing.T) {
	handler, mockGroupService, r := setupGroupTest()

	regularUserID := uuid.New()
	groupID := uuid.New()
	newMemberID := uuid.New()

	r.POST("/groups/:id/members", mockAuthMiddleware(regularUserID), handler.AddMember)

	// Mock service returning forbidden error
	forbiddenErr := errors.New("only admins can add members")
	mockGroupService.On("AddMember", mock.AnythingOfType("*context.timerCtx"), regularUserID, groupID, newMemberID).Return(forbiddenErr)

	// Prepare request
	body := map[string]string{"user_id": newMemberID.String()}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/groups/"+groupID.String()+"/members", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert - Should be 403 Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)

	mockGroupService.AssertExpectations(t)
}

func TestAddMember_BadRequest_InvalidGroupID(t *testing.T) {
	handler, _, r := setupGroupTest()

	adminID := uuid.New()
	r.POST("/groups/:id/members", mockAuthMiddleware(adminID), handler.AddMember)

	body := `{"user_id": "` + uuid.New().String() + `"}`
	req, _ := http.NewRequest("POST", "/groups/invalid-uuid/members", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddMember_BadRequest_MissingUserID(t *testing.T) {
	handler, _, r := setupGroupTest()

	adminID := uuid.New()
	groupID := uuid.New()

	r.POST("/groups/:id/members", mockAuthMiddleware(adminID), handler.AddMember)

	// Missing user_id
	body := `{}`
	req, _ := http.NewRequest("POST", "/groups/"+groupID.String()+"/members", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
