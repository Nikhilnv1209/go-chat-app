package handlers

import (
	"log"
	"net/http"

	"chat-app/internal/service"
	"chat-app/internal/websocket"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

type WSHandler struct {
	hub         *websocket.Hub
	authService service.AuthService
}

func NewWSHandler(hub *websocket.Hub, authService service.AuthService) *WSHandler {
	return &WSHandler{
		hub:         hub,
		authService: authService,
	}
}

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Helper to check origin for CORS
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for MVP. Lock down in production.
	},
}

func (h *WSHandler) ServeWS(c *gin.Context) {
	// 1. Auth Check (Token in Query Param)
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return
	}

	userID, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 2. Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	// 3. Register Client
	client := &websocket.Client{
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	h.hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}
