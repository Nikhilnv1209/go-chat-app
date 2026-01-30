package handlers

import (
	"context"
	"net/http"
	"time"

	"chat-app/internal/middleware"
	"chat-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupHandler struct {
	groupService service.GroupService
}

func NewGroupHandler(groupService service.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

type CreateGroupRequest struct {
	Name      string      `json:"name" binding:"required"`
	MemberIDs []uuid.UUID `json:"member_ids"`
}

type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

// CreateGroup handles POST /groups
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	// 1. Get user ID from AuthMiddleware context
	userID := middleware.GetUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. Parse request body
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// 3. Create group
	group, err := h.groupService.Create(ctx, userID, req.Name, req.MemberIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Return created group
	c.JSON(http.StatusCreated, gin.H{
		"id":   group.ID,
		"name": group.Name,
	})
}

// AddMember handles POST /groups/:id/members
func (h *GroupHandler) AddMember(c *gin.Context) {
	// 1. Get user ID from AuthMiddleware context
	adminID := middleware.GetUserIDFromContext(c)
	if adminID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. Parse group ID from URL
	groupIDStr := c.Param("id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	// 3. Parse request body
	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// 4. Add member
	err = h.groupService.AddMember(ctx, adminID, groupID, req.UserID)
	if err != nil {
		if err.Error() == "only admins can add members" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "member added successfully",
	})
}
