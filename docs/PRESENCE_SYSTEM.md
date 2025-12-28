# Presence System Documentation

## Overview

The presence system tracks and displays real-time online/offline status of users across the application. It uses a **hybrid approach** combining HTTP-based initial state with WebSocket-based real-time updates.

## Architecture

### Two-Phase Presence System

#### Phase 1: Initial State (HTTP)
- **When**: Page load, browser refresh, or conversation list refetch
- **How**: `GET /conversations` returns `is_online` field for DM conversations
- **Source**: PostgreSQL database (`users.is_online` column)
- **Purpose**: Provides accurate initial state without waiting for WebSocket

#### Phase 2: Real-Time Updates (WebSocket)
- **When**: Users connect/disconnect while you're online
- **How**: WebSocket events (`user_online`, `user_offline`)
- **Source**: In-memory Hub state (active connections)
- **Purpose**: Instant updates without polling or refreshing

---

## Backend Implementation

### Database Schema
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_online BOOLEAN DEFAULT FALSE,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### WebSocket Hub Flow

#### User Connection
```
Client WebSocket connects
    ↓
Hub.Register channel receives client
    ↓
Check: wasOffline = len(h.Clients[userID]) == 0
    ↓ (if TRUE - first connection)
    ├─→ updateUserStatus(userID, true)     → DB: SET is_online = TRUE
    └─→ broadcastPresence(userID, true)    → WS: Send "user_online" to contacts
    ↓
sendInitialPresence(client) → Send current online status of all contacts to THIS client
```

#### User Disconnection
```
Client WebSocket closes
    ↓
Hub.Unregister channel receives client
    ↓
Remove client from h.Clients[userID] array
    ↓
Check: isNowOffline = len(newClients) == 0
    ↓ (if TRUE - last connection closed)
    ├─→ updateUserStatus(userID, false)    → DB: SET is_online = FALSE
    └─→ broadcastPresence(userID, false)   → WS: Send "user_offline" to contacts
```

### HTTP API Response

**Endpoint**: `GET /conversations`

**Response Structure**:
```json
[
  {
    "id": "uuid",
    "type": "DM",
    "target_id": "user-uuid",
    "target_name": "alice",
    "last_message": "Hello!",
    "last_message_at": "2025-12-28T17:30:00Z",
    "unread_count": 3,
    "is_online": true     // ← Presence info
  },
  {
    "id": "uuid",
    "type": "GROUP",
    "target_id": "group-uuid",
    "target_name": "Team Chat",
    "last_message": "Meeting at 3",
    "last_message_at": "2025-12-28T17:25:00Z",
    "unread_count": 0,
    "member_count": 5
    // No is_online for groups
  }
]
```

---

## Frontend Implementation

### Redux State Management

**conversationSlice.ts** - Two key actions:

1. **`setConversations`**: Sets initial state from HTTP
```typescript
setConversations: (state, action: PayloadAction<Conversation[]>) => {
  // Backend is_online is the source of truth when fetched
  // Real-time WebSocket updates will override this via setUserOnlineStatus
  state.conversations = action.payload;
}
```

2. **`setUserOnlineStatus`**: Updates state from WebSocket
```typescript
setUserOnlineStatus: (state, action: PayloadAction<{ userId: string; isOnline: boolean }>) => {
  state.conversations.forEach((conv) => {
    if (conv.type === 'DM' && conv.target_id === action.payload.userId) {
      conv.is_online = action.payload.isOnline;
    }
  });
}
```

### WebSocket Hook

**useSocketConnection.ts** - Critical ordering:

```typescript
useEffect(() => {
  if (!isAuthenticated || !token || !user) {
    socketService.disconnect();
    return;
  }

  // STEP 1: Define and attach listeners FIRST
  const handleUserOnline = (payload: { user_id: string }) => {
    dispatch(setUserOnlineStatus({ userId: payload.user_id, isOnline: true }));
  };

  const handleUserOffline = (payload: { user_id: string }) => {
    dispatch(setUserOnlineStatus({ userId: payload.user_id, isOnline: false }));
  };

  socketService.on('user_online', handleUserOnline);
  socketService.on('user_offline', handleUserOffline);

  // STEP 2: THEN connect (ensures we don't miss events)
  socketService.connect(token);

  return () => {
    socketService.off('user_online', handleUserOnline);
    socketService.off('user_offline', handleUserOffline);
  };
}, [isAuthenticated, token, user, dispatch]);
```

**Why this order matters**: If we connect before attaching listeners, we might miss the `sendInitialPresence` events sent immediately after connection.

---

## Complete Presence Patterns

### Pattern 1: Fresh Page Load
```
Timeline:
1. User opens app (09:00 AM)
2. GET /conversations → Shows Alice: online, Bob: offline (from DB)
3. WebSocket connects
4. Backend sends user_online for Charlie (just came online)
5. UI updates: Charlie now online
```

**Result**: ✅ Accurate initial state + real-time updates

---

### Pattern 2: Browser Refresh
```
Timeline:
1. User refreshes page (09:05 AM)
2. Redux state cleared
3. GET /conversations → Fresh data from DB
4. WebSocket reconnects
5. Backend sends initial presence of all online contacts
```

**Result**: ✅ No flickering, accurate state restored

---

### Pattern 3: User Comes Online
```
Timeline:
1. Alice is viewing conversations
2. Bob opens app → WebSocket connects
3. Backend detects Bob's first connection (wasOffline = true)
4. Backend: updateUserStatus(Bob, true) → DB updated
5. Backend: broadcastPresence(Bob, true)
6. Alice receives: { type: "user_online", payload: { user_id: "Bob" } }
7. Redux: setUserOnlineStatus({ userId: "Bob", isOnline: true })
8. UI: Green dot appears next to Bob ✨
```

**Result**: ✅ Instant visual feedback

---

### Pattern 4: User Goes Offline
```
Timeline:
1. Alice is viewing conversations
2. Bob closes all tabs/browser
3. Bob's WebSocket closes → Hub.Unregister
4. Backend detects last connection closed (isNowOffline = true)
5. Backend: updateUserStatus(Bob, false) → DB updated
6. Backend: broadcastPresence(Bob, false)
7. Alice receives: { type: "user_offline", payload: { user_id: "Bob" } }
8. Redux: setUserOnlineStatus({ userId: "Bob", isOnline: false })
9. UI: Green dot disappears ✨
```

**Result**: ✅ Accurate offline detection

---

### Pattern 5: Multi-Tab User (Same Device)
```
Timeline:
1. Bob opens Tab 1 → WebSocket connects
   h.Clients[Bob] = [client1]
   wasOffline = true → Broadcasts "user_online"

2. Bob opens Tab 2 → WebSocket connects
   h.Clients[Bob] = [client1, client2]
   wasOffline = false → No broadcast ✅

3. Bob closes Tab 1
   h.Clients[Bob] = [client2]
   isNowOffline = false → No broadcast ✅

4. Bob closes Tab 2
   h.Clients[Bob] = []
   isNowOffline = true → Broadcasts "user_offline" ✅
```

**Result**: ✅ Correctly handles multiple connections per user

---

### Pattern 6: Network Disconnection (Ungraceful)
```
Timeline:
1. Bob's internet disconnects
2. WebSocket connection freezes (no proper close)
3. Backend's ReadPump doesn't receive pong
4. After pongWait timeout (60 seconds):
   - ReadPump exits
   - Defer triggers Hub.Unregister
   - User marked offline
```

**Result**: ⚠️ Delayed detection (up to 60s) - this is intentional to handle brief network glitches

---

### Pattern 7: React Query Refetch
```
Timeline:
1. User has conversations loaded, some users online
2. New message arrives → invalidateQueries(['conversations'])
3. GET /conversations refetches
4. setConversations(newData) called
5. New data includes is_online from DB
6. User's current online status preserved ✅
```

**Result**: ✅ No loss of presence state during refetches

---

## Edge Cases & Handling

### Edge Case 1: WebSocket Reconnection
**Scenario**: Connection drops and automatically reconnects

**Handling**:
```typescript
// socketService.ts
private attemptReconnect() {
  if (this.reconnectAttempts >= this.maxReconnectAttempts) {
    console.warn('Max reconnect attempts reached');
    return;
  }

  const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
  setTimeout(() => {
    this.reconnectAttempts++;
    if (this.token) {
      this.connect(this.token); // Reconnect with saved token
    }
  }, delay);
}
```

**On Reconnect**:
- Backend's `Register` channel receives new client
- `sendInitialPresence` re-sends current online contacts
- UI syncs back to accurate state

**Result**: ✅ Automatic recovery

---

### Edge Case 2: Race Between HTTP and WebSocket
**Scenario**: User loads page, WebSocket connects faster than HTTP response

**Timeline**:
```
T0: Page loads
T1: WebSocket connects (fast connection)
T2: Backend sends user_online events
T3: Listeners receive events → setUserOnlineStatus called
T4: HTTP GET /conversations returns
T5: setConversations called → Overwrites WS state
```

**Problem**: WebSocket updates lost when HTTP completes later

**Solution**: No longer an issue! Backend now provides `is_online` in HTTP response, so both sources provide accurate data.

---

### Edge Case 3: User Offline in DB But Actually Online
**Scenario**: Server crashed while user was connected, DB not updated

**Handling**:
1. User's next WebSocket connection triggers `wasOffline = true`
2. `updateUserStatus(user, true)` corrects DB
3. State synchronized ✅

**Result**: ✅ Self-healing system

---

### Edge Case 4: Listener Not Attached Before Initial Events
**Scenario**: `sendInitialPresence` fires before `useEffect` attaches listeners

**Previous Risk**:
```typescript
// BAD: Two separate effects
useEffect(() => { connect(token); }, [token]);           // Effect 1
useEffect(() => { on('user_online', handler); }, [...]); // Effect 2 - might run after!
```

**Solution**:
```typescript
// GOOD: Single effect with proper ordering
useEffect(() => {
  // 1. Attach listeners FIRST
  socketService.on('user_online', handleUserOnline);
  socketService.on('user_offline', handleUserOffline);

  // 2. THEN connect
  socketService.connect(token);

  return () => { /* cleanup */ };
}, [...]);
```

**Result**: ✅ Guaranteed to receive all events

---

## Visual Indicators

### ChatSidebar (Conversation List)
```tsx
{conv.is_online && conv.type === 'DM' && (
  <div className="absolute bottom-0 right-0 w-3.5 h-3.5 bg-green-500 border-2 border-white rounded-full"></div>
)}
```

### Chat Header (Active Conversation)
```tsx
<div className={cn("w-2 h-2 rounded-full",
  displayConversation.is_online ? "bg-green-500" : "bg-[#202022]/30"
)}></div>
<span className="text-xs text-[#202022]/50">
  {displayConversation.is_online ? 'Online' : 'Offline'}
</span>
```

---

## Performance Considerations

### Backend
- ✅ **O(1) lookups**: `h.Clients[userID]` is a map
- ✅ **Targeted broadcasts**: Only sends to users who have DM conversations with the changing user
- ✅ **No polling**: WebSocket push model
- ⚠️ **Trade-off**: 60-second timeout for ungraceful disconnects (prevents false positives from network blips)

### Frontend
- ✅ **Single source of truth**: Redux state
- ✅ **Optimistic updates**: WebSocket events update immediately
- ✅ **No redundant fetches**: React Query cache prevents unnecessary HTTP calls
- ✅ **Efficient rendering**: Only affected conversations re-render

---

## Testing Checklist

- [ ] Open app → See correct online status
- [ ] Refresh browser → Status remains accurate
- [ ] User comes online → Green dot appears
- [ ] User goes offline → Green dot disappears
- [ ] Open multiple tabs → No duplicate broadcasts
- [ ] Close one tab (still has another) → Stays online
- [ ] Close all tabs → Goes offline
- [ ] Disconnect WiFi → User goes offline after 60s
- [ ] Reconnect WiFi → User comes back online
- [ ] New conversation created → Invalidates list, refetches with presence

---

## Debugging

### Backend Logs
```go
log.Printf("User %s connected (wasOffline: %v)", userID, wasOffline)
log.Printf("User %s disconnected (isNowOffline: %v)", userID, isNowOffline)
log.Printf("Broadcasting presence: %s (%s) to %d contacts", userID, eventType, len(contacts))
```

### Frontend Logs
```
WS: Attaching event listeners...
WS: Connecting to socket...
WebSocket Connected
WS: User Online: <user_id>
WS: User Offline: <user_id>
```

### Common Issues

**Issue**: User shows offline but is actually online
- **Check**: Browser console for WS connection status
- **Check**: Backend logs for `wasOffline` flag
- **Fix**: Verify DB `is_online` column matches in-memory Hub state

**Issue**: Presence updates not received
- **Check**: Event listeners attached before connection?
- **Check**: WebSocket connection open?
- **Fix**: Review `useSocketConnection` hook ordering

**Issue**: Flashing offline/online on page load
- **Check**: Backend returning `is_online` in `GET /conversations`?
- **Fix**: Verify `ConversationResponse` includes `IsOnline` field

---

## Future Enhancements

1. **Last Seen Timestamp**: Show "Last seen 5 minutes ago" for offline users
2. **Typing Indicators**: Expand presence to include activity states
3. **Custom Status**: "Busy", "In a meeting", etc.
4. **Presence in Groups**: Show online member count in group chats
5. **Mobile Push**: Notify when specific users come online

---

## Summary

The hybrid presence system provides:
- ✅ **Accurate initial state** via HTTP
- ✅ **Real-time updates** via WebSocket
- ✅ **Multi-tab support** via connection counting
- ✅ **Graceful recovery** via reconnection logic
- ✅ **No race conditions** via proper event ordering
- ✅ **Efficient broadcasting** via targeted messaging

**Key Principle**: Trust the database for initial state, trust the WebSocket for live changes.
