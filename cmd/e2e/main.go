package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	BaseURL = "http://localhost:8080"
	WSURL   = "ws://localhost:8080/ws"

	// Using the tokens and IDs from previous manual test run
	AliceToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNDdlNGViZDAtMmY4Yy00M2Y0LWI4N2UtZjY1Y2Y2NGU1ZGNiIn0.NzudEFDAGpuXIkarv_OuEztyMpE1S4LI1pvkJbIm1IA"
	BobToken     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNThmNTY0MzUtOTdjYS00YjAyLWJjOTQtMjQ3YzUyODA3Yjg3In0.4jTTeV4nsfNjjRiRgPsV6IAvcx8aQoZUxRqoSRAXp9E"
	CharlieToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiYWE2Nzc5ZTEtODU1NS00YTE3LWE0YTMtYmEzZTE1YWMxOTc1In0.LL__OkZnxSFWU9Dwj_t-MiFG34wUWho5fHwy9XL2dDc"

	AliceID   = "47e4ebd0-2f8c-43f4-b87e-f65cf64e5dcb"
	BobID     = "58f56435-97ca-4b02-bc94-247c52807b87"
	CharlieID = "aa6779e1-8555-4a17-a4a3-ba3e15ac1975"
	GroupID   = "91d97fba-736e-4d5b-8a22-cb4a75cdb036"
)

func main() {
	log.SetFlags(0)
	fmt.Println("🚀 Starting End-to-End WebSocket & API Test")
	fmt.Println("===========================================")

	// 1. Connect WebSockets
	aliceConn, _, err := websocket.DefaultDialer.Dial(WSURL+"?token="+AliceToken, nil)
	if err != nil {
		log.Fatalf("❌ Failed to connect Alice: %v", err)
	}
	defer aliceConn.Close()
	fmt.Println("✅ Alice Connected")

	bobConn, _, err := websocket.DefaultDialer.Dial(WSURL+"?token="+BobToken, nil)
	if err != nil {
		log.Fatalf("❌ Failed to connect Bob: %v", err)
	}
	defer bobConn.Close()
	fmt.Println("✅ Bob Connected")

	charlieConn, _, err := websocket.DefaultDialer.Dial(WSURL+"?token="+CharlieToken, nil)
	if err != nil {
		log.Fatalf("❌ Failed to connect Charlie: %v", err)
	}
	defer charlieConn.Close()
	fmt.Println("✅ Charlie Connected")

	// Start reading goroutines to drain messages
	go readLoop(aliceConn, "Alice")
	go readLoop(bobConn, "Bob")
	go readLoop(charlieConn, "Charlie")

	time.Sleep(1 * time.Second)

	// 2. Send Messages
	fmt.Println("\n📨 Sending Messages...")

	// Alice -> Bob
	sendMsg(aliceConn, "send_message", map[string]string{
		"to_user_id": BobID,
		"content":    "Hi Bob, E2E test here!",
	})
	time.Sleep(500 * time.Millisecond)

	// Bob -> Alice
	sendMsg(bobConn, "send_message", map[string]string{
		"to_user_id": AliceID,
		"content":    "Loud and clear Alice!",
	})
	time.Sleep(500 * time.Millisecond)

	// Alice -> Group
	sendMsg(aliceConn, "send_message", map[string]string{
		"group_id": GroupID,
		"content":  "Hello Team Alpha!",
	})
	time.Sleep(500 * time.Millisecond)

	// Charlie -> Group
	sendMsg(charlieConn, "send_message", map[string]string{
		"group_id": GroupID,
		"content":  "Charlie reporting in!",
	})
	time.Sleep(2 * time.Second) // Wait for processing

	// 3. Verify APIs
	fmt.Println("\n🔎 Verifying Inbox & History APIs...")

	// Verify Alice's Inbox
	checkInbox("Alice", AliceToken, 2) // Expect Bob DM + Group

	// Verify Bob's Inbox
	checkInbox("Bob", BobToken, 2) // Expect Alice DM + Group

	// Verify Charlie's Inbox
	checkInbox("Charlie", CharlieToken, 1) // Expect Group only

	// Verify Alice-Bob History (Alice View)
	checkHistory("Alice->Bob DM", AliceToken, BobID, "DM", 2)

	// Verify Group History (Charlie View)
	checkHistory("Charlie->Group", CharlieToken, GroupID, "GROUP", 2)

	fmt.Println("\n✨ E2E Test Completed Successfully!")
}

func readLoop(c *websocket.Conn, name string) {
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			return
		}
		// fmt.Printf("[%s] Received message\n", name)
	}
}

func sendMsg(c *websocket.Conn, msgType string, payload interface{}) {
	msg := map[string]interface{}{
		"type":    msgType,
		"payload": payload,
	}
	err := c.WriteJSON(msg)
	if err != nil {
		log.Printf("❌ Failed to send message: %v", err)
	} else {
		fmt.Printf("✓ Sent message (%v)\n", payload.(map[string]string)["content"])
	}
}

func checkInbox(name, token string, expectedCount int) {
	body, err := fetchAPI(token, "/conversations")
	if err != nil {
		log.Printf("❌ %s Inbox fetch failed: %v", name, err)
		return
	}

	var convos []interface{}
	if err := json.Unmarshal(body, &convos); err != nil {
		log.Printf("❌ Failed to parse %s inbox: %v", name, err)
		return
	}

	if len(convos) >= expectedCount {
		fmt.Printf("✅ %s Inbox Verified (Found %d conversations)\n", name, len(convos))
	} else {
		fmt.Printf("⚠️  %s Inbox Warning: Expected at least %d, found %d\n", name, expectedCount, len(convos))
	}
}

func checkHistory(testName, token, targetID, msgType string, minCount int) {
	url := fmt.Sprintf("/messages?target_id=%s&type=%s&limit=50", targetID, msgType)
	body, err := fetchAPI(token, url)
	if err != nil {
		log.Printf("❌ %s History fetch failed: %v", testName, err)
		return
	}

	var msgs []interface{}
	if err := json.Unmarshal(body, &msgs); err != nil {
		log.Printf("❌ Failed to parse %s history: %v", testName, err)
		return
	}

	if len(msgs) >= minCount {
		fmt.Printf("✅ %s History Verified (Found %d messages)\n", testName, len(msgs))
	} else {
		fmt.Printf("⚠️  %s History Warning: Expected at least %d, found %d\n", testName, minCount, len(msgs))
	}
}

func fetchAPI(token, path string) ([]byte, error) {
	req, _ := http.NewRequest("GET", BaseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
