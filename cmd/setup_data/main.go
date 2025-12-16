package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const baseURL = "http://localhost:8080"
const wsURL = "ws://localhost:8080/ws"

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user"`
}

type GroupResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	fmt.Println("🚀 Starting Manual Test & Data Seeding...")
	fmt.Println("----------------------------------------")

	// 1. Register & Login Alice
	fmt.Println("👤 Setting up Alice...")
	aliceToken, aliceID := setupUser("Alice", "alice@test.com", "alice123")
	fmt.Printf("   -> Alice ID: %s\n", aliceID)

	// 2. Register & Login Bob
	fmt.Println("👤 Setting up Bob...")
	bobToken, bobID := setupUser("Bob", "bob@test.com", "bob123")
	fmt.Printf("   -> Bob ID: %s\n", bobID)

	// 3. Connect WebSockets
	fmt.Println("\n🔌 Connecting to WebSockets...")
	aliceWS, _ := connectWS(aliceToken, "Alice")
	bobWS, _ := connectWS(bobToken, "Bob")
	defer aliceWS.Close()
	defer bobWS.Close()

	// 4. Send DM: Alice -> Bob
	fmt.Println("\n💬 Testing DM sending (Alice -> Bob)...")
	sendDM(aliceToken, bobID, "Hey Bob! This is a test DM from Alice.")

	// Wait a bit for WS delivery
	time.Sleep(500 * time.Millisecond)

	// 5. Create Group
	fmt.Println("\n👥 Testing Group Creation...")
	groupID := createGroup(aliceToken, "Engineering Team", "Work discussions")
	fmt.Printf("   -> Group created: %s\n", groupID)

	// 6. Add Bob to Group
	fmt.Println("\n➕ Adding Bob to Group...")
	addMember(aliceToken, groupID, bobID)

	// 7. Send Group Messages
	fmt.Println("\n📢 Testing Group Messaging...")
	sendGroupMessage(aliceToken, groupID, "Welcome to the team, Bob!")
	time.Sleep(500 * time.Millisecond)
	sendGroupMessage(bobToken, groupID, "Thanks Alice! Happy to be here.")

	fmt.Println("\n----------------------------------------")
	fmt.Println("✅ Data Seeding & Logic Check Complete!")
	fmt.Println("   You can now reload the frontend to see the conversations.")
}

// --- Helper Functions ---

func setupUser(username, email, password string) (string, string) {
	// Register
	client := &http.Client{}
	regPayload := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}
	jsonData, _ := json.Marshal(regPayload)

	// Try Register (ignore error if exists)
	resp, err := client.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err == nil {
		resp.Body.Close()
	}

	// Login
	loginResp, err := client.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to login %s: %v", username, err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != 200 {
		log.Fatalf("Login failed for %s: Status %d", username, loginResp.StatusCode)
	}

	var res LoginResponse
	body, _ := io.ReadAll(loginResp.Body)
	json.Unmarshal(body, &res)

	return res.Token, res.User.ID
}

func sendDM(token, targetID, content string) {
	// client := &http.Client{}

	// Create or Get Conversation first (implicit in some systems, checks implementation)
	// Actually, most chat apps just send message and conversation is created if not exists.
	// But let's cheat and just call the WebSocket send route? No, the user asked to check routes.
	// We don't have a direct REST content creation route in the specs usually,
	// typically it's over WS or specific endpoint.
	// Checking the `chat_handler.go`, we only have GET routes for history.
	// Messages are usually sent via WebSocket in this app architecture (based on spec).
	// Let's check if there is a POST /messages endpoint.
	// Looking at `main.go`, there isn't a POST /messages route!
	// Sending messages is PURELY WebSocket based in this implementation?
	// Wait, let me check `service.MessageService`.

	// Checking main.go again...
	// chatRoutes.GET("/conversations", chatHandler.GetConversations)
	// chatRoutes.GET("/messages", chatHandler.GetMessages)
	// The ONLY POST routes are for Auth and Groups.
	// AND POST /messages/:id/read

	// THIS MEANS MESSAGE SENDING IS ONLY VIA WEBSOCKET.
	// So I must use the WebSocket to send the messages.

	// Let's use the WS connection we established.
	fmt.Printf("   -> Sending via WS: %s\n", content)

	msg := map[string]interface{}{
		"type": "send_message",
		"payload": map[string]interface{}{
			"recipient_id": targetID,
			"content":      content,
			"type":         "DM",
		},
	}

	sendWSMessage(token, msg)
}

func sendGroupMessage(token, groupID, content string) {
	fmt.Printf("   -> Sending Group Msg via WS: %s\n", content)
	msg := map[string]interface{}{
		"type": "send_message",
		"payload": map[string]interface{}{
			"group_id": groupID,
			"content":  content,
			"type":     "GROUP",
		},
	}
	sendWSMessage(token, msg)
}

// Global map to hold connections for the script
var conns = make(map[string]*websocket.Conn)

func connectWS(token, name string) (*websocket.Conn, error) {
	u, _ := url.Parse(wsURL)
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("❌ %s failed to connect to WS: %v", name, err)
		return nil, err
	}
	fmt.Printf("   ✅ %s connected to WS\n", name)

	// Start reading pump to consume messages (receipts/echoes)
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			// Just verify we got something
			fmt.Printf("      [WS %s] Received: %s\n", name, string(message))
		}
	}()

	conns[token] = c
	return c, nil
}

func sendWSMessage(token string, msg interface{}) {
	conn := conns[token]
	if conn == nil {
		fmt.Println("   ❌ No WS connection for this token to send message")
		return
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("   ❌ Failed to write JSON to WS: %v", err)
	}
}

func createGroup(token, name, description string) string {
	client := &http.Client{}
	payload := map[string]string{
		"name":        name,
		"description": description,
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", baseURL+"/groups", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to create group: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Group creation failed status %d: %s", resp.StatusCode, string(body))
	}

	var res GroupResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &res)
	return res.ID
}

func addMember(token, groupID, userID string) {
	client := &http.Client{}
	payload := map[string]string{
		"user_id": userID,
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/groups/%s/members", baseURL, groupID), bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to add member: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("   ⚠️ Add member warning (might already be member): %s\n", string(body))
	}
}
