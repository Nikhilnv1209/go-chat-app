package websocket

import (
	"context"
	"chat-app/internal/service"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessagePayload struct {
	ToUserID uuid.UUID `json:"to_user_id"` // Simplified for DM
	GroupID  uuid.UUID `json:"group_id"`   // For Group (optional)
	Content  string    `json:"content"`
}

type MessageDeliveredPayload struct {
	MessageID uuid.UUID `json:"message_id"`
}

type TypingPayload struct {
	ConversationType string    `json:"conversation_type"` // "DM" or "GROUP"
	TargetID         uuid.UUID `json:"target_id"`         // user ID for DM, group ID for GROUP
}

type SetActiveConversationPayload struct {
	ConversationType string    `json:"conversation_type"` // "DM" or "GROUP"
	TargetID         uuid.UUID `json:"target_id"`         // user ID for DM, group ID for GROUP
}

// HandleMessage routes incoming WS messages to appropriate services
func HandleMessage(message []byte, client *Client, msgService service.MessageService) {
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return
	}

	switch wsMsg.Type {
	case "set_active_conversation":
		var payload SetActiveConversationPayload
		if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
			log.Printf("Invalid Payload for set_active_conversation: %v", err)
			return
		}

		// Validate conversation type
		if payload.ConversationType != "DM" && payload.ConversationType != "GROUP" {
			log.Printf("Invalid conversation type: %s", payload.ConversationType)
			return
		}

		// Set or clear active conversation
		if payload.TargetID == uuid.Nil {
			client.Hub.ClearActiveConversation(client)
		} else {
			client.Hub.SetActiveConversation(client, payload.ConversationType, payload.TargetID)
		}

	case "send_message":
		var payload SendMessagePayload
		if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
			log.Printf("Invalid Payload: %v", err)
			return
		}

		ctx := context.Background()
		if payload.ToUserID != uuid.Nil {
			// Direct Message
			msg, err := msgService.SendDirectMessage(ctx, client.UserID, payload.ToUserID, payload.Content)
			if err != nil {
				log.Printf("Failed to send DM: %v", err)
				return
			}

			// Ack to Sender
			ack, _ := json.Marshal(map[string]interface{}{
				"type":    "message_sent",
				"payload": msg,
			})
			client.Send <- ack
		} else if payload.GroupID != uuid.Nil {
			// Group Message
			msg, err := msgService.SendGroupMessage(ctx, client.UserID, payload.GroupID, payload.Content)
			if err != nil {
				log.Printf("Failed to send group message: %v", err)
				return
			}

			// Ack to Sender
			ack, _ := json.Marshal(map[string]interface{}{
				"type":    "message_sent",
				"payload": msg,
			})
			client.Send <- ack
		}

	case "message_delivered":
		var payload MessageDeliveredPayload
		if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
			log.Printf("Invalid Payload for message_delivered: %v", err)
			return
		}

		ctx := context.Background()
		if err := msgService.MarkAsDelivered(ctx, client.UserID, []uuid.UUID{payload.MessageID}); err != nil {
			log.Printf("Failed to mark delivered: %v", err)
		}

	case "typing_start":
		handleTypingStart(client, wsMsg.Payload, msgService)

	case "typing_stop":
		handleTypingStop(client, wsMsg.Payload, msgService)

	default:
		log.Printf("Unknown message type: %s", wsMsg.Type)
	}
}

// handleTypingStart broadcasts typing indicator to relevant users
func handleTypingStart(client *Client, payload json.RawMessage, msgService service.MessageService) {
	var typingPayload TypingPayload
	if err := json.Unmarshal(payload, &typingPayload); err != nil {
		log.Printf("Invalid typing_start payload: %v", err)
		return
	}

	// Validation
	if typingPayload.ConversationType != "DM" && typingPayload.ConversationType != "GROUP" {
		log.Printf("Invalid conversation type: %s", typingPayload.ConversationType)
		return
	}
	if typingPayload.TargetID == uuid.Nil {
		log.Printf("Invalid target_id for typing event")
		return
	}

	ctx := context.Background()
	// Get user info for the typing user
	user, err := msgService.GetUserInfo(ctx, client.UserID)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		return
	}
	if user == nil {
		log.Printf("User info not found for ID: %s", client.UserID)
		return
	}

	// Broadcast typing event
	if err := msgService.BroadcastTypingIndicator(ctx, client.UserID, user.Username, typingPayload.ConversationType, typingPayload.TargetID, true); err != nil {
		log.Printf("Failed to broadcast typing_start: %v", err)
	}
}

// handleTypingStop broadcasts typing stop indicator to relevant users
func handleTypingStop(client *Client, payload json.RawMessage, msgService service.MessageService) {
	var typingPayload TypingPayload
	if err := json.Unmarshal(payload, &typingPayload); err != nil {
		log.Printf("Invalid typing_stop payload: %v", err)
		return
	}

	// Validation
	if typingPayload.ConversationType != "DM" && typingPayload.ConversationType != "GROUP" {
		return // Silent ignore for stop events to handle spam/race conditions gracefully
	}

	ctx := context.Background()
	// Broadcast typing stop event
	if err := msgService.BroadcastTypingIndicator(ctx, client.UserID, "", typingPayload.ConversationType, typingPayload.TargetID, false); err != nil {
		log.Printf("Failed to broadcast typing_stop: %v", err)
	}
}
