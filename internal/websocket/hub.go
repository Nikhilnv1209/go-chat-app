package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"chat-app/internal/repository"

	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients (UserID -> List of Clients)
	Clients map[uuid.UUID][]*Client

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Lock for safely updating the Clients map
	mu sync.RWMutex

	// Repository to update user status
	userRepo repository.UserRepository

	// Repository to find user's contacts for presence broadcast
	convRepo repository.ConversationRepository

	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewHub(userRepo repository.UserRepository, convRepo repository.ConversationRepository) *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uuid.UUID][]*Client),
		userRepo:   userRepo,
		convRepo:   convRepo,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.ctx.Done():
			h.shutdown()
			return

		case client := <-h.Register:
			h.mu.Lock()
			// If first connection, mark online and broadcast
			wasOffline := len(h.Clients[client.UserID]) == 0
			h.Clients[client.UserID] = append(h.Clients[client.UserID], client)
			h.mu.Unlock()

			if wasOffline {
				h.wg.Add(1)
				go h.updateUserStatus(client.UserID, true)
				h.wg.Add(1)
				go h.broadcastPresence(client.UserID, true)
			}

			// Send initial presence of contacts to the new client
			h.wg.Add(1)
			go h.sendInitialPresence(client)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Clients[client.UserID]; ok {
				// Filter out the client being unregistered
				newClients := make([]*Client, 0, len(clients)-1)
				for _, c := range clients {
					if c != client {
						newClients = append(newClients, c)
					}
				}

				isNowOffline := len(newClients) == 0
				if isNowOffline {
					delete(h.Clients, client.UserID)
				} else {
					h.Clients[client.UserID] = newClients
				}
				close(client.Send)
				h.mu.Unlock()

				if isNowOffline {
					h.wg.Add(1)
					go h.updateUserStatus(client.UserID, false)
					h.wg.Add(1)
					go h.broadcastPresence(client.UserID, false)
				}
			} else {
				h.mu.Unlock()
			}
		}
	}
}

func (h *Hub) updateUserStatus(userID uuid.UUID, isOnline bool) {
	defer h.wg.Done()

	// Create context with timeout for DB operation
	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()

	if err := h.userRepo.UpdateOnlineStatus(ctx, userID, isOnline, time.Now()); err != nil {
		// Log error but don't crash
		log.Printf("Failed to update user status for %s: %v", userID, err)
	}
}

// Shutdown gracefully closes the hub
func (h *Hub) Shutdown(ctx context.Context) error {
	log.Println("Hub: Initiating shutdown...")
	h.cancel()

	done := make(chan struct{})
	go func() {
		h.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Hub: All goroutines finished cleanly")
		return nil
	case <-ctx.Done():
		log.Println("Hub: Shutdown timed out, forcing exit")
		return ctx.Err()
	}
}

// shutdown performs internal cleanup
func (h *Hub) shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for userID, clients := range h.Clients {
		for _, client := range clients {
			close(client.Send)
		}
		delete(h.Clients, userID)
		log.Printf("Hub: Closed %d connections for user %s", len(clients), userID)
	}

	log.Println("Hub: Shutdown complete")
}

// broadcastPresence sends presence updates to all users who have a DM conversation with this user
func (h *Hub) broadcastPresence(userID uuid.UUID, isOnline bool) {
	defer h.wg.Done()

	// Create context with timeout for DB operation
	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()

	eventType := "user_offline"
	if isOnline {
		eventType = "user_online"
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"type": eventType,
		"payload": map[string]interface{}{
			"user_id": userID,
		},
	})

	// Find all users who have this user as a target_id in their DM conversations
	// These are the users that need to know about this user's presence
	contacts, err := h.convRepo.FindContactsOfUser(ctx, userID)
	if err != nil {
		log.Printf("Failed to find contacts for presence broadcast: %v", err)
		return
	}

	for _, contactID := range contacts {
		h.SendToUser(contactID, payload)
	}
}

// SendToUser sends a message to all connected devices of a specific user.
func (h *Hub) SendToUser(userID uuid.UUID, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.Clients[userID]; ok {
		for _, client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				// Clean up locked client in next loop or let ReadPump handle it
			}
		}
	}
}

// IsUserViewingConversation checks if any client of the user is currently viewing the specified conversation.
// Returns true if at least one client has this conversation as active.
func (h *Hub) IsUserViewingConversation(convType string, targetID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	targetConv := convType + ":" + targetID.String()

	// Check if any client of the target user has this conversation active
	// Note: This is called from message sending context where we check the RECEIVER
	// We need to find all clients and check if any of them is viewing this conversation
	for _, clients := range h.Clients {
		for _, client := range clients {
			if client.ActiveConversation == targetConv {
				return true
			}
		}
	}

	return false
}

// SetActiveConversation sets the active conversation for a specific client.
func (h *Hub) SetActiveConversation(client *Client, convType string, targetID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client.ActiveConversation = convType + ":" + targetID.String()
}

// ClearActiveConversation clears the active conversation for a specific client.
func (h *Hub) ClearActiveConversation(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client.ActiveConversation = ""
}

// sendInitialPresence sends the current online status of contacts to a newly connected client
func (h *Hub) sendInitialPresence(client *Client) {
	defer h.wg.Done()

	// Create context with timeout for DB operation
	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()

	// 1. Get user's conversations
	convs, err := h.convRepo.FindByUser(ctx, client.UserID)
	if err != nil {
		log.Printf("Failed to fetch conversations for initial presence sync: %v", err)
		return
	}

	// 2. Filter for DMs and checking online status
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conv := range convs {
		if conv.Type == "DM" {
			targetID := conv.TargetID
			// Check if target is online (has active clients)
			if clients, ok := h.Clients[targetID]; ok && len(clients) > 0 {
				// 3. Send 'user_online' event to THIS client only
				payload, _ := json.Marshal(map[string]interface{}{
					"type": "user_online",
					"payload": map[string]interface{}{
						"user_id": targetID,
					},
				})
				select {
				case client.Send <- payload:
				case <-h.ctx.Done():
					return
				}
			}
		}
	}
}
