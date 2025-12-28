package handlers

import (
	"net/http"
	"time"

	"chat-app/internal/errors"
	"chat-app/internal/middleware"
	"chat-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrValidation.Message, "details": err.Error()})
		return
	}

	accessToken, refreshToken, user, err := h.service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.setRefreshTokenCookie(c, refreshToken)

	c.JSON(http.StatusCreated, gin.H{
		"token": accessToken,
		"user":  user,
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrValidation.Message, "details": err.Error()})
		return
	}

	accessToken, refreshToken, user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.setRefreshTokenCookie(c, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
		"user":  user,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	accessToken, newRefreshToken, err := h.service.Refresh(refreshToken)
	if err != nil {
		// If refresh fails, clear cookie
		h.clearRefreshTokenCookie(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Update cookie if it changed (optional rotation)
	if newRefreshToken != "" && newRefreshToken != refreshToken {
		h.setRefreshTokenCookie(c, newRefreshToken)
	}

	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		_ = h.service.Logout(refreshToken)
	}

	h.clearRefreshTokenCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) SearchUsers(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	query := c.Query("q")
	users, err := h.service.SearchUsers(query, userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *AuthHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Status, gin.H{
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}

// Helpers

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	// Determine if we're in production (Secure flag should be true for HTTPS)
	isProduction := gin.Mode() == gin.ReleaseMode

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		MaxAge:   7 * 24 * 60 * 60,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction, // true in production (HTTPS), false in development (HTTP)
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	isProduction := gin.Mode() == gin.ReleaseMode

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}
