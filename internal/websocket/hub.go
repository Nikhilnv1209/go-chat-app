package websocket

import (
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
}

func NewHub(userRepo repository.UserRepository) *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uuid.UUID][]*Client),
		userRepo:   userRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			// If first connection, mark online
			if len(h.Clients[client.UserID]) == 0 {
				go h.updateUserStatus(client.UserID, true)
			}
			h.Clients[client.UserID] = append(h.Clients[client.UserID], client)
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Clients[client.UserID]; ok {
				// Filter out the client being unregistered
				newClients := make([]*Client, 0)
				for _, c := range clients {
					if c != client {
						newClients = append(newClients, c)
					}
				}

				if len(newClients) == 0 {
					delete(h.Clients, client.UserID)
					// Verify this was the last connection before marking offline
					// (Double check inside lock isn't strictly necessary if strict ordering is guaranteed,
					// but good for safety if we add async stuff)
					go h.updateUserStatus(client.UserID, false)
				} else {
					h.Clients[client.UserID] = newClients
				}
				close(client.Send)
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) updateUserStatus(userID uuid.UUID, isOnline bool) {
	// We might want to pass context later
	if err := h.userRepo.UpdateOnlineStatus(userID, isOnline, time.Now()); err != nil {
		// Log error but don't crash
		log.Printf("Failed to update user status for %s: %v", userID, err)
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
