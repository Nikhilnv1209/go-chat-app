package handlers

import (
	"net/http"
	"strconv"

	"chat-app/internal/middleware"
	"chat-app/internal/repository"
	"chat-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	convRepo   repository.ConversationRepository
	msgRepo    repository.MessageRepository
	userRepo   repository.UserRepository
	groupRepo  repository.GroupRepository
	msgService service.MessageService
}

func NewChatHandler(
	convRepo repository.ConversationRepository,
	msgRepo repository.MessageRepository,
	userRepo repository.UserRepository,
	groupRepo repository.GroupRepository,
	msgService service.MessageService,
) *ChatHandler {
	return &ChatHandler{
		convRepo:   convRepo,
		msgRepo:    msgRepo,
		userRepo:   userRepo,
		groupRepo:  groupRepo,
		msgService: msgService,
	}
}

// ConversationResponse is the response structure for conversation list
type ConversationResponse struct {
	ID            uuid.UUID `json:"id"`
	Type          string    `json:"type"`
	TargetID      uuid.UUID `json:"target_id"`
	TargetName    string    `json:"target_name"`
	LastMessageAt string    `json:"last_message_at"`
	UnreadCount   int       `json:"unread_count"`
}

// GetConversations handles GET /conversations
// Returns the inbox (list of conversations) sorted by last_message_at
func (h *ChatHandler) GetConversations(c *gin.Context) {
	// 1. Get user ID from AuthMiddleware context
	userID := middleware.GetUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. Fetch conversations for this user
	conversations, err := h.convRepo.FindByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch conversations"})
		return
	}

	// 3. Build response with target names
	response := make([]ConversationResponse, 0, len(conversations))
	for _, conv := range conversations {
		targetName := ""

		switch conv.Type {
		case "DM":
			// Fetch user name
			user, err := h.userRepo.FindByID(conv.TargetID)
			if err == nil {
				targetName = user.Username
			}
		case "GROUP":
			// Fetch group name
			group, err := h.groupRepo.FindByID(conv.TargetID)
			if err == nil {
				targetName = group.Name
			}
		}

		response = append(response, ConversationResponse{
			ID:            conv.ID,
			Type:          conv.Type,
			TargetID:      conv.TargetID,
			TargetName:    targetName,
			LastMessageAt: conv.LastMessageAt.Format("2006-01-02T15:04:05Z07:00"),
			UnreadCount:   conv.UnreadCount,
		})
	}

	// 4. Return response
	c.JSON(http.StatusOK, response)
}

// GetMessages handles GET /messages?target_id=<uuid>&type=<DM|GROUP>&limit=<n>
// Returns message history for a specific conversation
func (h *ChatHandler) GetMessages(c *gin.Context) {
	// 1. Get user ID from AuthMiddleware context
	userID := middleware.GetUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. Parse query parameters
	targetIDStr := c.Query("target_id")
	if targetIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_id is required"})
		return
	}

	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target_id"})
		return
	}

	msgType := c.Query("type")
	if msgType == "" {
		msgType = "DM" // Default to DM
	}
	if msgType != "DM" && msgType != "GROUP" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be DM or GROUP"})
		return
	}

	limit := 50 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Parse 'before_id' cursor for pagination (message UUID)
	var beforeID *uuid.UUID
	if beforeIDStr := c.Query("before_id"); beforeIDStr != "" {
		parsed, err := uuid.Parse(beforeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'before_id' format, must be a valid UUID"})
			return
		}
		beforeID = &parsed
	}

	// 3. Verify user has access to this conversation
	if msgType == "GROUP" {
		// Check if user is a member of the group
		isMember, err := h.groupRepo.IsMember(targetID, userID)
		if err != nil || !isMember {
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not a member of this group"})
			return
		}
	}

	// 4. Fetch messages
	messages, err := h.msgService.GetHistory(userID, targetID, msgType, limit, beforeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	// 5. Reset unread count for this conversation
	_ = h.convRepo.ResetUnread(userID, msgType, targetID)

	// 6. Return messages
	c.JSON(http.StatusOK, messages)
}

// MarkRead handles POST /messages/:id/read
func (h *ChatHandler) MarkRead(c *gin.Context) {
	// 1. Get user ID from AuthMiddleware context
	userID := middleware.GetUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. Parse ID param
	messageIDStr := c.Param("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		return
	}

	// 3. Mark as read
	if err := h.msgService.MarkAsRead(userID, []uuid.UUID{messageID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "READ", "message_id": messageID})
}

func (h *ChatHandler) GetReceipts(c *gin.Context) {
	// 1. Get user ID from AuthMiddleware context
	userID := middleware.GetUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	messageIDStr := c.Param("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	receipts, err := h.msgService.GetMessageReceipts(userID, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receipts)
}
