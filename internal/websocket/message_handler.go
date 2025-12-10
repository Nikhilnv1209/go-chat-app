package websocket

import (
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

// HandleMessage routes incoming WS messages to appropriate services
func HandleMessage(message []byte, client *Client, msgService service.MessageService) {
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return
	}

	switch wsMsg.Type {
	case "send_message":
		var payload SendMessagePayload
		if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
			log.Printf("Invalid Payload: %v", err)
			return
		}

		if payload.ToUserID != uuid.Nil {
			// Direct Message
			msg, err := msgService.SendDirectMessage(client.UserID, payload.ToUserID, payload.Content)
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
			msg, err := msgService.SendGroupMessage(client.UserID, payload.GroupID, payload.Content)
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

		if err := msgService.MarkAsDelivered(client.UserID, []uuid.UUID{payload.MessageID}); err != nil {
			log.Printf("Failed to mark delivered: %v", err)
		}

	default:
		log.Printf("Unknown message type: %s", wsMsg.Type)
	}
}
