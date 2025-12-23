# Technical Blueprint

**Status**: Draft
**Version**: 1.0
**Owner**: System Architect Agent

## 1. Architecture Overview
A **Client-Server** architecture where Next.js serves as both the secure frontend host (RSC) and the BFF (Backend-for-Frontend) if needed, though we connect primarily directly to the Go API.

`Browser` <-> `Next.js (App Router)` <-> `Go API (REST + WS)`

## 2. Directory Structure (Next.js 16 - IMPLEMENTED)

```bash
frontend/
├── app/
│   ├── (auth)/             # Route Group: Public
│   │   ├── layout.tsx      # Auth layout with background
│   │   ├── login/page.tsx  # ✅ COMPLETED
│   │   ├── register/page.tsx # ✅ COMPLETED
│   ├── dashboard/          # Route Group: Protected (changed from (dashboard))
│   │   ├── layout.tsx      # ✅ COMPLETED: NavigationRail + ChatSidebar + AuthGuard + Responsive
│   │   ├── page.tsx        # ✅ COMPLETED: "Select a chat" placeholder
│   │   ├── profile/        # ✅ BONUS: Profile page
│   │   │   └── page.tsx
│   ├── layout.tsx          # ✅ COMPLETED: Root with providers
│   ├── page.tsx            # Landing page
│   ├── providers.tsx       # ✅ COMPLETED: Redux/Query providers
│   └── globals.css         # ✅ COMPLETED: Tailwind styles
├── components/
│   ├── ui/                 # ✅ COMPLETED: Shadcn components
│   │   ├── button.tsx
│   │   ├── input.tsx
│   │   ├── card.tsx
│   │   ├── avatar.tsx
│   │   ├── badge.tsx
│   │   ├── scroll-area.tsx
│   │   ├── textarea.tsx
│   │   └── tooltip.tsx     # ✅ ADDED
│   ├── chat/               # Chat Organisms
│   │   ├── ChatSidebar.tsx # ✅ COMPLETED: Conversation list
│   │   └── UserProfile.tsx # ✅ COMPLETED: User info component
│   └── dashboard/          # Dashboard Components
│       └── NavigationRail.tsx # ✅ ADDED: Modern sidebar navigation
├── lib/
│   ├── api.ts              # ✅ COMPLETED: API wrapper with Axios
│   ├── conversationApi.ts  # ✅ BONUS: Conversation-specific API
│   └── utils.ts            # ✅ COMPLETED: cn() utility
├── store/                  # ✅ COMPLETED: Redux Toolkit
│   ├── store.ts            # Root store configuration
│   ├── hooks.ts            # Redux hooks
│   └── features/
│       ├── authSlice.ts    # ✅ COMPLETED: Auth state management
│       └── conversationSlice.ts # ✅ BONUS: Conversation state
├── types/                  # ✅ BONUS: Type definitions
│   └── index.ts            # All TypeScript interfaces
└── package.json            # ✅ COMPLETED: Dependencies and scripts
```

**NOTES ON IMPLEMENTATION**:
- Changed route group from `(dashboard)` to `dashboard` (removed parentheses)
- Added bonus features: Profile page, conversationSlice, and separate conversationApi
- WebSocket service (`socket.ts`) not yet implemented
- Message components not yet created

## 3. Data Strategy

### 3.1 Authentication
*   **Storage**: JWT stored in `localStorage` (MVP) or HTTP-Only Cookies (Production).
*   **Auth Guard**: A `useAuthType` hook in the `(dashboard)` layout checks for token presence. Redirects to `/login` if missing.

### 3.2 State Management Diagram

| Data Type | Owner | Persistence | Sync Mechanism |
|---|---|---|---|
| **User Session** | Redux (`authSlice`) | LocalStorage | Login Response |
| **Theme / Sidebar** | Redux (`uiSlice`) | LocalStorage | UI Toggle |
| **Inbox List** | React Query | Cache (5min) | REST `GET /conversations` |
| **Messages (History)** | React Query | Cache (Infinite) | REST `GET /messages` |
| **New Messages** | React Query | Manual Cache Update | WebSocket Event |

### 3.3 WebSocket Integration
The `SocketService` is a singleton class exposed via `useSocket`.

**Event Flow**:
1.  Connection established on `(dashboard)` mount.
2.  `onMessage` event received.
3.  `queryClient.setQueryData(['messages', id], (old) => [...old, newMessage])`
    *   *Note*: This bypasses Redux to keep message streams performant.

## 4. Dependencies (INSTALLED)
*   `next`: 16.0.8 (Latest stable)
*   `react`: 19.2.1 (Stable)
*   `@tanstack/react-query`: 5.90.12
*   `@reduxjs/toolkit`: 2.11.1
*   `react-redux`: 9.2.0
*   `lucide-react`: 0.560.0 (Icon library)
*   `axios`: 1.13.2 (HTTP client)
*   `shadcn/ui` components (CLI-based component library)
    *   Uses Radix UI primitives for accessibility
    *   Components in `@/components/ui`
    *   Full customization, zero runtime overhead
    *   Installed: `@radix-ui/react-avatar`, `@radix-ui/react-scroll-area`, `@radix-ui/react-slot`
*   `clsx`: 2.1.1 (For dynamic classes)
*   `tailwind-merge`: 3.4.0 (Tailwind class merging)
*   `class-variance-authority`: 0.7.1 (For component variants)
*   `tailwindcss`: v4 (Latest version)

## 5. Error Handling Strategy

### Client-Side Errors
| Error Type | Handling | User Feedback |
|------------|----------|---------------|
| **Network Failure** | React Query `onError` callback | Toast: "Connection lost. Retrying..." |
| **API 4xx Errors** | Catch in mutation/query | Inline error message on forms |
| **API 5xx Errors** | Global error boundary | Full-page error with "Retry" button |
| **WebSocket Disconnect** | Exponential backoff reconnect | Toast: "Reconnecting..." |

### Error Boundary Hierarchy
```
<RootErrorBoundary> (500 page)
  └── <QueryErrorResetBoundary>
        └── <DashboardErrorBoundary> (Per-feature fallback)
```

## 6. Security Considerations

| Threat | Mitigation |
|--------|------------|
| **XSS (Cross-Site Scripting)** | React's default escaping + never use `dangerouslySetInnerHTML` with user content |
| **CSRF (Cross-Site Request Forgery)** | JWT in Authorization header (not cookies) + SameSite cookie policy if using cookies |
| **Token Theft** | Store JWT in memory (sessionStorage) for MVP; upgrade to HTTP-only cookies for production |
| **Insecure WebSocket** | Validate token on every WS connection; reject invalid tokens server-side |
| **Input Validation** | Client-side validation with Zod; never trust client data on server |

### Secure Defaults
*   All external links open with `rel="noopener noreferrer"`.
*   Content Security Policy (CSP) headers recommended in production.

## 7. API Contract Reference

The frontend consumes the following backend endpoints:

| Endpoint | Method | Purpose | Spec Reference |
|----------|--------|---------|----------------|
| `/auth/register` | POST | Create new user account | Backend F01 |
| `/auth/login` | POST | Authenticate and receive JWT | Backend F01 |
| `/conversations` | GET | Fetch user's inbox (sorted by recent) | Backend F05 |
| `/messages` | GET | Fetch message history | Backend F05 |
| `/messages/:id/read` | POST | Mark message as read | Backend F06 |
| `/messages/:id/receipts` | GET | Get receipt status for a message | Backend F06 |
| `/groups` | POST | Create a new group | Backend F04 |
| `/groups/:id/members` | POST | Add member to group (admin only) | Backend F04 |
| `/ws` | WS | Real-time messaging | Backend F02/F03/F06/F07 |

*Full API documentation: See `README.md` > API Documentation section.*

### 7.1 Error Response Format

All API errors return a consistent JSON structure:

```json
{
  "error": "descriptive message"
}
```

| Status Code | Meaning | Common Causes |
|-------------|---------|---------------|
| `400` | Bad Request | Invalid input, missing required fields |
| `401` | Unauthorized | Missing/invalid/expired JWT token |
| `403` | Forbidden | Access denied (e.g., not group member) |
| `404` | Not Found | Resource doesn't exist |
| `500` | Server Error | Internal backend error |

### 7.2 Pagination Strategy

**Implementation**: ID-based cursor pagination using message UUIDs.

```
Initial load:     GET /messages?target_id=<uuid>&type=DM&limit=50
Next page:        GET /messages?target_id=<uuid>&type=DM&limit=50&before_id=<oldest-msg-id>
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `target_id` | UUID | Required | User ID (DM) or Group ID (GROUP) |
| `type` | String | `DM` | Conversation type: `DM` or `GROUP` |
| `limit` | Integer | `50` | Max messages to return |
| `before_id` | UUID | Optional | Cursor: fetch messages older than this message ID |

**Frontend Strategy for Infinite Scroll**:
1. Initial load: `GET /messages?target_id=<id>&type=DM&limit=50`
2. Store the oldest message ID from response
3. On scroll up: `GET /messages?...&before_id=<oldest_msg_id>`
4. Repeat until response is empty (no more history)

## 8. WebSocket Event Schema (Frontend Handling)

### 8.1 Connection

```javascript
const ws = new WebSocket(`ws://localhost:8080/ws?token=${jwt}`);
```

**Reconnection Strategy**:
- Implement exponential backoff: 1s → 2s → 4s → 8s → 16s → max 30s
- Refresh JWT if expired before reconnecting
- Maximum 10 retry attempts before showing error UI

### 8.2 Outgoing Events (Client → Server)

| Event Type | Payload Schema | When to Send |
|------------|----------------|--------------|
| `send_message` | `{ to_user_id?: string, group_id?: string, content: string }` | User sends a message |
| `typing_start` | `{ conversation_type: "DM" \| "GROUP", target_id: string }` | User begins typing (debounce 300ms) |
| `typing_stop` | `{ conversation_type: "DM" \| "GROUP", target_id: string }` | 3s after last keystroke OR message sent |
| `message_delivered` | `{ message_id: string }` | After displaying received message |

**Example: Send Message**
```json
{
  "type": "send_message",
  "payload": {
    "to_user_id": "uuid-of-recipient",
    "content": "Hello, World!"
  }
}
```

**Example: Typing Start (Group)**
```json
{
  "type": "typing_start",
  "payload": {
    "conversation_type": "GROUP",
    "target_id": "uuid-of-group"
  }
}
```

### 8.3 Incoming Events (Server → Client)

| Event Type | Payload Schema | Frontend Action |
|------------|----------------|-----------------|
| `new_message` | `Message` object | Add to message list, update conversation order, increment unread |
| `message_sent` | `Message` object | Confirm optimistic UI, replace temp ID with real ID |
| `user_typing` | `{ user_id, username, conversation_type, target_id }` | Show "Alice is typing..." indicator |
| `user_stopped_typing` | `{ user_id, conversation_type, target_id }` | Hide typing indicator |
| `receipt_update` | `{ message_id, user_id, status, updated_at }` | Update message status icon (✓ → ✓✓ → blue ✓✓) |

**Example: New Message**
```json
{
  "type": "new_message",
  "payload": {
    "id": "msg-uuid",
    "sender_id": "sender-uuid",
    "receiver_id": "your-uuid",
    "content": "Hello!",
    "msg_type": "private",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

**Example: Message Sent (Acknowledgment)**
```json
{
  "type": "message_sent",
  "payload": {
    "id": "msg-uuid-from-server",
    "sender_id": "your-uuid",
    "receiver_id": "recipient-uuid",
    "content": "Hello!",
    "msg_type": "private",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

**Example: User Typing**
```json
{
  "type": "user_typing",
  "payload": {
    "user_id": "uuid-of-typing-user",
    "username": "Alice",
    "conversation_type": "DM",
    "target_id": "conversation-uuid"
  }
}
```

**Example: Receipt Update**
```json
{
  "type": "receipt_update",
  "payload": {
    "message_id": "msg-uuid",
    "user_id": "reader-uuid",
    "status": "READ",
    "updated_at": "2024-01-01T12:01:00Z"
  }
}
```

### 8.4 Receipt Status Icons

| Status | Icon | Meaning |
|--------|------|---------|
| `SENT` | ✓ (gray) | Message saved on server |
| `DELIVERED` | ✓✓ (gray) | Message delivered to recipient's device |
| `READ` | ✓✓ (blue) | Recipient opened the conversation |

## 9. TypeScript Type Definitions

```typescript
// types/api.ts

export interface User {
  id: string;
  username: string;
  email: string;
  is_online: boolean;
  last_seen: string;
  created_at: string;
}

export interface Message {
  id: string;
  sender_id: string;
  receiver_id: string | null;
  group_id: string | null;
  content: string;
  msg_type: 'private' | 'group';
  created_at: string;
}

export interface Conversation {
  id: string;
  type: 'DM' | 'GROUP';
  target_id: string;
  target_name: string;
  last_message_at: string;
  unread_count: number;
}

export interface MessageReceipt {
  id: string;
  message_id: string;
  user_id: string;
  status: 'SENT' | 'DELIVERED' | 'READ';
  created_at: string;
  updated_at: string;
}

export interface Group {
  id: string;
  name: string;
  created_at: string;
}

// WebSocket Event Types
export type WSOutgoingEvent =
  | { type: 'send_message'; payload: { to_user_id?: string; group_id?: string; content: string } }
  | { type: 'typing_start'; payload: { conversation_type: 'DM' | 'GROUP'; target_id: string } }
  | { type: 'typing_stop'; payload: { conversation_type: 'DM' | 'GROUP'; target_id: string } }
  | { type: 'message_delivered'; payload: { message_id: string } };

export type WSIncomingEvent =
  | { type: 'new_message'; payload: Message }
  | { type: 'message_sent'; payload: Message }
  | { type: 'user_typing'; payload: { user_id: string; username: string; conversation_type: string; target_id: string } }
  | { type: 'user_stopped_typing'; payload: { user_id: string; conversation_type: string; target_id: string } }
  | { type: 'receipt_update'; payload: { message_id: string; user_id: string; status: string; updated_at: string } };
```
