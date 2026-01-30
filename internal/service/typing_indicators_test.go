package service_test

import (
	"context"
	"chat-app/internal/models"
	"chat-app/internal/service"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Tests for F07 - Typing Indicators

func TestBroadcastTypingIndicator_DM_TypingStart(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)
	mockHub := new(MockHub)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	senderID := uuid.New()
	targetID := uuid.New()
	senderUsername := "Alice"

	// Mock: Hub.SendToUser should be called with typing event
	mockHub.On("SendToUser", targetID, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)

		if msg["type"] != "user_typing" {
			return false
		}

		payloadData := msg["payload"].(map[string]interface{})
		return payloadData["user_id"] == senderID.String() &&
			payloadData["username"] == senderUsername &&
			payloadData["conversation_type"] == "DM" &&
			payloadData["target_id"] == targetID.String()
	})).Return()

	// Execute
	err := svc.BroadcastTypingIndicator(ctx, senderID, senderUsername, "DM", targetID, true)

	// Assert
	assert.NoError(t, err)
	mockHub.AssertExpectations(t)
}

func TestBroadcastTypingIndicator_DM_TypingStop(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)
	mockHub := new(MockHub)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	senderID := uuid.New()
	targetID := uuid.New()

	// Mock: Hub.SendToUser should be called with typing stop event
	mockHub.On("SendToUser", targetID, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)

		if msg["type"] != "user_stopped_typing" {
			return false
		}

		payloadData := msg["payload"].(map[string]interface{})
		// Username should NOT be present for typing stop
		_, hasUsername := payloadData["username"]
		return !hasUsername &&
			payloadData["user_id"] == senderID.String() &&
			payloadData["conversation_type"] == "DM" &&
			payloadData["target_id"] == targetID.String()
	})).Return()

	// Execute
	err := svc.BroadcastTypingIndicator(ctx, senderID, "", "DM", targetID, false)

	// Assert
	assert.NoError(t, err)
	mockHub.AssertExpectations(t)
}

func TestBroadcastTypingIndicator_Group_Success(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)
	mockHub := new(MockHub)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	senderID := uuid.New()
	groupID := uuid.New()
	member1 := uuid.New()
	member2 := uuid.New()
	senderUsername := "Carol"

	// Mock: Sender is a member of the group
	mockGroupRepo.On("IsMember", ctx, groupID, senderID).Return(true, nil)

	// Mock: Get group members
	members := []models.GroupMember{
		{GroupID: groupID, UserID: senderID, Role: "MEMBER"},
		{GroupID: groupID, UserID: member1, Role: "MEMBER"},
		{GroupID: groupID, UserID: member2, Role: "MEMBER"},
	}
	mockGroupRepo.On("GetMembers", ctx, groupID).Return(members, nil)

	// Mock: Hub.SendToUser should be called for EACH OTHER member (not sender)
	mockHub.On("SendToUser", member1, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)
		return msg["type"] == "user_typing"
	})).Return()

	mockHub.On("SendToUser", member2, mock.MatchedBy(func(payload []byte) bool {
		var msg map[string]interface{}
		json.Unmarshal(payload, &msg)
		return msg["type"] == "user_typing"
	})).Return()

	// Execute
	err := svc.BroadcastTypingIndicator(ctx, senderID, senderUsername, "GROUP", groupID, true)

	// Assert
	assert.NoError(t, err)
	mockHub.AssertCalled(t, "SendToUser", member1, mock.Anything)
	mockHub.AssertCalled(t, "SendToUser", member2, mock.Anything)
	// Sender should NOT receive the typing event
	mockHub.AssertNotCalled(t, "SendToUser", senderID, mock.Anything)
	mockGroupRepo.AssertExpectations(t)
}

func TestBroadcastTypingIndicator_Group_NotMember(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)
	mockHub := new(MockHub)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	nonMemberID := uuid.New()
	groupID := uuid.New()

	// Mock: Sender is NOT a member
	mockGroupRepo.On("IsMember", ctx, groupID, nonMemberID).Return(false, nil)

	// Execute
	err := svc.BroadcastTypingIndicator(ctx, nonMemberID, "Hacker", "GROUP", groupID, true)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "sender is not a member of the group", err.Error())
	mockGroupRepo.AssertExpectations(t)
	// Hub should NOT be called
	mockHub.AssertNotCalled(t, "SendToUser", mock.Anything, mock.Anything)
}

func TestGetUserInfo_Success(t *testing.T) {
	ctx := context.Background()
	mockMsgRepo := new(MockMessageRepo)
	mockConvRepo := new(MockConversationRepo)
	mockGroupRepo := new(MockGroupRepo)
	mockReceiptRepo := new(MockMessageReceiptRepo)
	mockUserRepo := new(MockUserRepo)
	mockHub := new(MockHub)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo, mockGroupRepo, mockReceiptRepo, mockUserRepo, mockHub)

	userID := uuid.New()
	expectedUser := &models.User{
		BaseModel: models.BaseModel{ID: userID},
		Username:  "TestUser",
		Email:     "test@example.com",
	}

	// Mock
	mockUserRepo.On("FindByID", ctx, userID).Return(expectedUser, nil)

	// Execute
	user, err := svc.GetUserInfo(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "TestUser", user.Username)
	mockUserRepo.AssertExpectations(t)
}
