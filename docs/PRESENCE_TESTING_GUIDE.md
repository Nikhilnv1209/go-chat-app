# Presence System - Complete Manual Testing Guide

## Test Setup

### Prerequisites
1. **Two test users** in the database:
   - User A: `alice@test.com` / `alice123`
   - User B: `bob@test.com` / `bob123`
2. **DM conversation** between Alice and Bob (send at least one message)
3. **Two browsers** or browser profiles (Chrome normal + Chrome incognito, or Chrome + Firefox)
4. **Backend running**: `make dev` or `cd backend && go run cmd/server/main.go`
5. **Frontend running**: `npm run dev`

### Browser Setup
- **Browser A**: Will be logged in as Alice
- **Browser B**: Will be logged in as Bob
- **DevTools Open**: Press F12 in both browsers to monitor Console and Network tabs

---

## Test Suite

### ✅ TEST 1: Basic Initial State

**Objective**: Verify presence loads correctly on first app open

#### Steps:
1. **Close all browser tabs** for both users completely
2. **Wait 65 seconds** (pong timeout + buffer) to ensure both users are marked offline in DB
3. Open **Browser A** → Navigate to `http://localhost:3000`
4. Login as Alice
5. **OBSERVE**: Conversation list loads

#### Expected Results:
- ✅ Bob's conversation appears in the list
- ✅ Bob shows **OFFLINE** (no green dot)
- ✅ Console shows: `GET /conversations` returns `is_online: false` for Bob
- ✅ No flickering or state changes

#### Browser Console Check:
```
Network Tab → GET /conversations → Preview:
[
  {
    "id": "...",
    "type": "DM",
    "target_name": "bob",
    "is_online": false  // ← Should be false
  }
]
```

---

### ✅ TEST 2: User Comes Online (Real-Time Update)

**Objective**: Verify WebSocket `user_online` event updates UI instantly

#### Steps:
1. **Keep Browser A open** (Alice still logged in)
2. Open **Browser B** → Navigate to `http://localhost:3000`
3. Login as Bob
4. **OBSERVE Browser A** (Alice's screen)

#### Expected Results:
- ✅ Bob's avatar/conversation **immediately** shows green dot
- ✅ If Alice is viewing Bob's chat, status changes to "Online"
- ✅ No page refresh needed

#### Browser A Console Check (Alice):
```
Console output:
WS: User Online: <bob-user-id>
```

#### Browser B Console Check (Bob):
```
Console output:
WebSocket Connected
```

---

### ✅ TEST 3: User Goes Offline (Real-Time Update)

**Objective**: Verify WebSocket `user_offline` event works

#### Steps:
1. **Both browsers open** (Alice and Bob logged in)
2. Close **Browser B** completely (Bob's browser)
3. **OBSERVE Browser A** (Alice's screen)

#### Expected Results:
- ✅ Within **1-2 seconds**, Bob's green dot disappears
- ✅ Status changes to "Offline"
- ✅ Change happens without Alice refreshing

#### Browser A Console Check (Alice):
```
Console output:
WS: User Offline: <bob-user-id>
```

---

### ✅ TEST 4: Browser Refresh (Both Users Online)

**Objective**: Verify presence persists after refresh, no flash of offline state

#### Steps:
1. **Both browsers open** (Alice and Bob logged in)
2. Verify Bob shows as **ONLINE** in Alice's browser
3. In **Browser A**, press `Ctrl+R` or `F5` to refresh
4. **OBSERVE carefully** during the reload

#### Expected Results:
- ✅ After refresh, Bob **immediately** shows as ONLINE
- ✅ **NO flickering** from offline → online
- ✅ Green dot visible as soon as conversations render

#### Browser A Console Check (Alice):
```
Network Tab → GET /conversations:
{
  "target_name": "bob",
  "is_online": true  // ← Should be true from DB
}

Console:
WS: Connecting to socket...
WebSocket Connected
WS: User Online: <bob-user-id>  // ← Sent from sendInitialPresence
```

---

### ✅ TEST 5: Browser Refresh (Other User Offline)

**Objective**: Verify offline state also persists correctly on refresh

#### Steps:
1. **Close Browser B** completely (Bob logs out)
2. Wait **5 seconds** for Bob to be marked offline
3. Verify Bob shows as **OFFLINE** in Browser A
4. In **Browser A**, refresh the page (`Ctrl+R`)

#### Expected Results:
- ✅ After refresh, Bob shows as **OFFLINE** immediately
- ✅ No green dot visible
- ✅ No flash of online state

---

### ✅ TEST 6: Multi-Tab - Same User

**Objective**: Verify user stays online when opening multiple tabs

#### Steps:
1. **Browser A**: Alice logged in (should show Bob as offline)
2. In **same browser**, open **new tab** → Login as Bob (`Tab B1`)
3. **OBSERVE Browser A**: Bob should now show online
4. In **same browser as Bob**, open **another new tab** → Login as Bob again (`Tab B2`)
5. **OBSERVE Browser A**: Bob should still show online (no duplicate events)
6. Close `Tab B1` (Bob's first tab)
7. **OBSERVE Browser A**: Bob should **REMAIN online** (Tab B2 still open)
8. Close `Tab B2` (Bob's last tab)
9. **OBSERVE Browser A**: Bob should now go **OFFLINE**

#### Expected Results:
- ✅ Opening 2nd tab: No duplicate `user_online` event
- ✅ Closing 1st tab: Bob stays online
- ✅ Closing last tab: Bob goes offline
- ✅ Backend correctly tracks multiple connections per user

#### Backend Logs Check:
```
Step 3: User <bob-id> connected (wasOffline: true)
        Broadcasting presence: user_online to X contacts

Step 4: User <bob-id> connected (wasOffline: false)
        [NO broadcast - not the first connection]

Step 6: User <bob-id> disconnected (isNowOffline: false)
        [NO broadcast - still has another connection]

Step 8: User <bob-id> disconnected (isNowOffline: true)
        Broadcasting presence: user_offline to X contacts
```

---

### ✅ TEST 7: Multi-Device - Same User

**Objective**: Same as TEST 6 but across different browsers/devices

#### Steps:
1. Open **Chrome** → Login as Bob
2. Open **Firefox** (or Chrome Incognito) → Login as Bob
3. Alice should see Bob as **ONLINE**
4. Close **Chrome** (Bob's first device)
5. Alice should see Bob **STILL ONLINE** (Firefox still open)
6. Close **Firefox** (Bob's last device)
7. Alice should see Bob go **OFFLINE**

#### Expected Results:
- Same as TEST 6

---

### ✅ TEST 8: Hard Refresh (Ctrl+Shift+R)

**Objective**: Verify presence survives cache clearing

#### Steps:
1. Both users logged in, Bob online
2. In Browser A, press `Ctrl+Shift+R` (hard refresh, clears cache)
3. **OBSERVE**

#### Expected Results:
- ✅ After hard refresh, Bob still shows as ONLINE
- ✅ No flickering or incorrect state
- ✅ WebSocket reconnects successfully

---

### ✅ TEST 9: Network Disconnection (Ungraceful)

**Objective**: Verify ungraceful disconnect is detected

#### Steps:
1. Both users logged in
2. In Browser B (Bob):
   - Open DevTools → Network tab
   - Click "Offline" to simulate network loss
3. Wait **65 seconds** (ping/pong timeout)
4. **OBSERVE Browser A**

#### Expected Results:
- ✅ After ~60-65 seconds, Bob goes OFFLINE in Alice's view
- ✅ This delay is intentional (prevents false positives from brief network glitches)

#### Alternative Test (Faster):
1. In Browser B, open DevTools → Console
2. Run: `performance.navigation.reload()`
3. Immediately observe Browser A
4. Bob should briefly go offline, then back online when reconnection completes

---

### ✅ TEST 10: WebSocket Reconnection

**Objective**: Verify automatic reconnection works and syncs state

#### Steps:
1. Both users logged in, Bob online
2. In Browser A (Alice):
   - Open DevTools → Console
   - Type: `socketService.disconnect()`
   - Press Enter
3. Wait **3 seconds**
4. Check Browser A console

#### Expected Results:
- ✅ Console shows: `WebSocket Closed`
- ✅ Console shows: `Attempting reconnect in <delay>ms...`
- ✅ Console shows: `Connecting to WebSocket:`
- ✅ Console shows: `WebSocket Connected`
- ✅ Bob's online status **re-syncs** and shows correctly
- ✅ Alice receives `user_online` event for Bob from `sendInitialPresence`

---

### ✅ TEST 11: Backend Restart (Both Users Connected)

**Objective**: Verify state recovery after backend crash/restart

#### Steps:
1. Both users logged in, Bob shows as online
2. **Stop backend**: `Ctrl+C` in backend terminal
3. **OBSERVE both browsers**: WebSocket connections close
4. **Restart backend**: `make dev`
5. Wait for backend to fully start
6. **OBSERVE both browsers**

#### Expected Results:
- ✅ Both browsers automatically reconnect
- ✅ Console shows reconnection attempts
- ✅ After reconnect:
  - Database `is_online` fields are fresh
  - Both users re-establish connections
  - Presence syncs correctly

#### Check Database:
```bash
# Before restart (both online)
SELECT username, is_online FROM users;
# alice | true
# bob   | true

# After restart (both reconnected)
SELECT username, is_online FROM users;
# alice | true
# bob   | true
```

---

### ✅ TEST 12: Race Between HTTP and WebSocket

**Objective**: Ensure no race condition when page loads

#### Steps:
1. Both users logged in
2. In Browser A (Alice):
   - Open DevTools → Network tab
   - Click "Slow 3G" throttling
3. Refresh page (`Ctrl+R`)
4. **OBSERVE loading sequence**

#### Expected Results:
- ✅ WebSocket may connect before HTTP response
- ✅ Alice may receive `user_online` for Bob via WebSocket first
- ✅ HTTP response arrives later with `is_online: true`
- ✅ **No conflict**: Both sources say Bob is online
- ✅ Final state is correct: Bob shows as ONLINE

---

### ✅ TEST 13: New Conversation (No Prior Messages)

**Objective**: Verify presence works for newly created conversations

#### Steps:
1. Create a third user: `charlie@test.com` / `charlie123`
2. Login as Charlie in a new browser tab
3. Alice (Browser A) should NOT have a conversation with Charlie yet
4. In Browser A (Alice):
   - Click "New Chat"
   - Search for "charlie"
   - Start a conversation, send a message
5. **OBSERVE**: Charlie's conversation appears in sidebar

#### Expected Results:
- ✅ New conversation shows Charlie as **ONLINE** (green dot)
- ✅ Status is accurate from the start
- ✅ WebSocket `conversation_created` event received
- ✅ List invalidates and refetches with presence data

---

### ✅ TEST 14: Presence in Chat Header

**Objective**: Verify presence shows in active chat view

#### Steps:
1. Both users logged in, Bob online
2. In Browser A, click on Bob's conversation
3. **OBSERVE chat header**
4. Close Browser B (Bob goes offline)
5. **OBSERVE chat header updates**

#### Expected Results:
- ✅ When Bob online: Green dot + "Online" text
- ✅ When Bob offline: Gray dot + "Offline" text
- ✅ Updates in real-time without leaving chat

---

### ✅ TEST 15: Presence in Sidebar

**Objective**: Verify presence shows in conversation list

#### Steps:
1. Both users logged in
2. In Browser A, stay on "All Chats" view (don't open specific chat)
3. Verify Bob's conversation has green dot on avatar
4. Close Browser B
5. Green dot should disappear

#### Expected Results:
- ✅ Green dot appears on bottom-right of avatar when online
- ✅ Disappears when offline
- ✅ Updates without clicking into the chat

---

### ✅ TEST 16: Logout/Login (Same Browser)

**Objective**: Verify presence updates on authentication changes

#### Steps:
1. Both users logged in, Bob online in Browser A
2. In Browser A:
   - Click logout
   - **OBSERVE**: Redirected to login page
   - Login again as Alice
3. **OBSERVE** conversation list

#### Expected Results:
- ✅ On logout: WebSocket disconnects
- ✅ On login: Fresh fetch shows Bob's correct status
- ✅ WebSocket reconnects
- ✅ Presence is accurate

---

### ✅ TEST 17: Multiple Conversations

**Objective**: Verify presence updates across all conversations correctly

#### Steps:
1. Create 3 users: Alice, Bob, Charlie
2. Alice has DM conversations with both Bob and Charlie
3. Bob and Charlie both login (both online)
4. **OBSERVE Browser A**: Both show green dots
5. Close Bob's browser
6. **OBSERVE Browser A**: Only Bob goes offline, Charlie stays online
7. Close Charlie's browser
8. **OBSERVE Browser A**: Both now offline

#### Expected Results:
- ✅ Each conversation tracked independently
- ✅ Status changes only affect the correct user
- ✅ No cross-contamination of presence state

---

### ✅ TEST 18: Group Conversations (NO Presence)

**Objective**: Verify groups don't show presence indicators

#### Steps:
1. Create a group with Alice, Bob, and Charlie
2. Login as Alice
3. **OBSERVE** the group conversation in sidebar

#### Expected Results:
- ✅ Group shows Users icon (not individual avatar)
- ✅ **NO green dot** on group (groups don't have is_online)
- ✅ Member count shown instead
- ✅ No errors in console

---

### ✅ TEST 19: Rapid Connection/Disconnection

**Objective**: Stress test - verify no duplicates or missed updates

#### Steps:
1. Alice logged in (Browser A)
2. Bob rapidly:
   - Open browser → Login (connect)
   - Close browser (disconnect)
   - Repeat 5 times quickly
3. **OBSERVE Browser A** and backend logs

#### Expected Results:
- ✅ Each connection triggers `user_online`
- ✅ Each disconnection triggers `user_offline`
- ✅ No duplicate broadcasts
- ✅ Final state matches reality (if Bob closed browser, shows offline)

---

### ✅ TEST 20: Long Session (Ping/Pong)

**Objective**: Verify connection stays alive with ping/pong

#### Steps:
1. Both users logged in
2. Leave both browsers **open and idle** for **5 minutes**
3. Don't interact with either browser
4. After 5 minutes, send a message from Alice to Bob
5. **OBSERVE**

#### Expected Results:
- ✅ Both WebSocket connections stay alive (ping/pong keeps them alive)
- ✅ Both users still show as online
- ✅ Message sends successfully
- ✅ No reconnection needed

#### Console Check:
```
# Should NOT see:
WebSocket Closed
Attempting reconnect...

# Should see (every ~54 seconds - pingPeriod):
[Backend sends ping frames automatically]
```

---

### ✅ TEST 21: Browser Tab Sleep (Chrome Power Saving)

**Objective**: Verify presence when browser throttles inactive tabs

#### Steps:
1. Both users logged in
2. In Browser B (Bob):
   - Switch to a **different tab** (not the chat app)
   - Leave it in background for **10 minutes**
3. In Browser A (Alice), check Bob's status
4. Switch back to Browser B (Bob's tab)

#### Expected Results:
- ✅ Bob may briefly go offline if Chrome throttled the tab
- ✅ On tab focus, reconnection should occur
- ✅ Bob comes back online
- ⚠️ This is a known browser behavior, acceptable delay

---

### ✅ TEST 22: Database Manual Update

**Objective**: Verify system self-heals from DB inconsistency

#### Steps:
1. Bob is online (Browser B open)
2. **Manually update database**:
```sql
UPDATE users SET is_online = false WHERE username = 'bob';
```
3. Alice refreshes her page (Browser A)
4. **OBSERVE**: Alice sees Bob as **OFFLINE** (wrong - he's actually online)
5. Bob closes browser
6. Bob opens browser again (reconnects)
7. **OBSERVE Browser A**

#### Expected Results:
- ✅ Step 7: Bob shows as **ONLINE** (DB corrected on connection)
- ✅ System self-healed from inconsistent state

---

### ✅ TEST 23: Concurrent Logins (Same User, Same Time)

**Objective**: Verify simultaneous connections handled correctly

#### Steps:
1. Alice logged in (Browser A)
2. Open **2 new browsers** simultaneously
3. In **both new browsers**, login as Bob at the **exact same time**
4. **OBSERVE Browser A** and backend logs

#### Expected Results:
- ✅ Only **1** `user_online` broadcast (first connection sets wasOffline=true)
- ✅ Second connection doesn't duplicate broadcast
- ✅ Alice sees Bob as online (correct)

---

### ✅ TEST 24: JWT Token Expiry

**Objective**: Verify presence on session timeout

#### Steps:
1. Both users logged in
2. Wait for access token to expire (default: 15 minutes)
3. Refresh **Browser A** (Alice)
4. **OBSERVE**

#### Expected Results:
- ✅ Token refresh happens automatically
- ✅ WebSocket reconnects with new token
- ✅ Presence re-syncs correctly
- ✅ Bob's status still accurate

---

### ✅ TEST 25: Backend Under Load

**Objective**: Verify presence works with high message traffic

#### Steps:
1. Both users logged in
2. In Browser B (Bob), rapidly send **50 messages** to Alice
3. While messages are flying, close Browser B
4. **OBSERVE Browser A**

#### Expected Results:
- ✅ Messages appear correctly
- ✅ Bob's status changes to **OFFLINE** despite message traffic
- ✅ No presence events lost in message flood

---

## Edge Case Tests

### 🔶 EDGE 1: Server Shutdown While User Connects

**Steps**:
1. Start backend
2. Open Browser A, start logging in (click submit)
3. **Immediately** stop backend (`Ctrl+C`)
4. **OBSERVE**

**Expected**:
- ✅ Login fails gracefully
- ✅ Frontend shows connection error
- ✅ No crash or infinite loading

---

### 🔶 EDGE 2: Invalid WebSocket Token

**Steps**:
1. Login as Alice
2. Open DevTools → Console
3. Run: `socketService.disconnect()`
4. Run: `socketService.connect('invalid-token-xyz')`

**Expected**:
- ✅ Backend rejects connection (401 or immediate close)
- ✅ Frontend shows console error
- ✅ Reconnection attempts with correct token

---

### 🔶 EDGE 3: User Deleted While Online

**Steps**:
1. Bob logged in
2. **Delete Bob from database**:
```sql
DELETE FROM users WHERE username = 'bob';
```
3. Bob tries to send a message
4. Alice checks Bob's presence

**Expected**:
- ✅ Bob's message fails (user not found)
- ✅ Alice may still see Bob as online (cached state)
- ✅ On next refresh, conversation might error or disappear

---

### 🔶 EDGE 4: Simultaneous Disconnect of All Users

**Steps**:
1. 10 users logged in
2. **Restart backend** (all disconnect simultaneously)
3. All users' browsers auto-reconnect

**Expected**:
- ✅ Backend handles reconnection flood
- ✅ All users re-establish connections
- ✅ Presence syncs for all

---

## Automated Testing (Developer Notes)

While manual testing is comprehensive, consider adding:

```bash
# Backend unit tests
go test ./internal/websocket/... -v

# Frontend integration tests (Cypress/Playwright)
describe('Presence System', () => {
  it('shows user online when they connect', () => {
    // ... automated version of TEST 2
  });
});
```

---

## Test Results Template

Use this checklist to track your testing:

```
[ ] TEST 1: Basic Initial State
[ ] TEST 2: User Comes Online
[ ] TEST 3: User Goes Offline
[ ] TEST 4: Browser Refresh (Both Online)
[ ] TEST 5: Browser Refresh (Other Offline)
[ ] TEST 6: Multi-Tab Same User
[ ] TEST 7: Multi-Device Same User
[ ] TEST 8: Hard Refresh
[ ] TEST 9: Network Disconnection
[ ] TEST 10: WebSocket Reconnection
[ ] TEST 11: Backend Restart
[ ] TEST 12: Race Between HTTP/WS
[ ] TEST 13: New Conversation
[ ] TEST 14: Presence in Chat Header
[ ] TEST 15: Presence in Sidebar
[ ] TEST 16: Logout/Login
[ ] TEST 17: Multiple Conversations
[ ] TEST 18: Group Conversations
[ ] TEST 19: Rapid Connect/Disconnect
[ ] TEST 20: Long Session
[ ] TEST 21: Browser Tab Sleep
[ ] TEST 22: Database Manual Update
[ ] TEST 23: Concurrent Logins
[ ] TEST 24: JWT Token Expiry
[ ] TEST 25: Backend Under Load

EDGE CASES:
[ ] EDGE 1: Server Shutdown During Login
[ ] EDGE 2: Invalid WebSocket Token
[ ] EDGE 3: User Deleted While Online
[ ] EDGE 4: Simultaneous Disconnect
```

---

## Success Criteria

**ALL tests should pass with:**
- ✅ No console errors
- ✅ Accurate presence state
- ✅ Real-time updates < 2 seconds
- ✅ No UI flickering
- ✅ Graceful error handling

If any test fails, check:
1. Backend logs
2. Frontend console
3. Network tab WebSocket frames
4. Database state (`SELECT * FROM users;`)

---

## Quick Debugging Commands

### Check Database State
```sql
SELECT username, is_online, last_seen FROM users;
```

### Check Backend Connections
```go
// Add to hub.go temporarily
func (h *Hub) PrintConnections() {
    h.mu.RLock()
    defer h.mu.RUnlock()
    for userID, clients := range h.Clients {
        log.Printf("User %s: %d connections", userID, len(clients))
    }
}
```

### Check Frontend WebSocket State
```javascript
// Browser Console
socketService.socket?.readyState
// 0 = CONNECTING, 1 = OPEN, 2 = CLOSING, 3 = CLOSED
```

---

**Good luck with testing! 🚀**
