# Technical Blueprint

**Status**: Draft
**Version**: 1.0
**Owner**: System Architect Agent

## 1. Architecture Overview
A **Client-Server** architecture where Next.js serves as both the secure frontend host (RSC) and the BFF (Backend-for-Frontend) if needed, though we connect primarily directly to the Go API.

`Browser` <-> `Next.js (App Router)` <-> `Go API (REST + WS)`

## 2. Directory Structure (Next.js 15)

```bash
frontend/
├── app/
│   ├── (auth)/             # Route Group: Public
│   │   ├── login/page.tsx
│   │   └── register/page.tsx
│   ├── (dashboard)/        # Route Group: Protected
│   │   ├── layout.tsx      # Sidebar + AuthGuard
│   │   ├── page.tsx        # "Select a chat" placeholder
│   │   └── c/[id]/page.tsx # Active Conversation
│   ├── layout.tsx          # Root: ThemeProvider, ReduxProvider, QueryProvider
│   └── globals.css
├── components/
│   ├── ui/                 # Shadcn Atoms (Button, Input, Avatar)
│   ├── auth/               # Auth Forms
│   └── chat/               # Chat Organisms
│       ├── ChatSidebar.tsx
│       ├── MessageList.tsx
│       ├── ChatInput.tsx
│       └── UserProfile.tsx
├── lib/
│   ├── api.ts              # Typed Fetch Wrapper
│   └── socket.ts           # Singleton WebSocket Service
├── store/                  # Redux Toolkit
│   ├── store.ts
│   └── features/
│       ├── authSlice.ts
│       └── uiSlice.ts
└── hooks/                  # Custom Hooks
    ├── useSocket.ts
    └── useChatScroll.ts
```

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

## 4. Dependencies
*   `next`: ^15.0.0
*   `react`: ^19.0.0 (RC/Stable)
*   `@tanstack/react-query`: ^5.0.0
*   `@reduxjs/toolkit`: ^2.0.0
*   `react-redux`: ^9.0.0
*   `lucide-react`: Latest
*   `clsx`, `tailwind-merge`: For dynamic classes.

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
| `/conversations` | GET | Fetch user's inbox | Backend F05 |
| `/messages` | GET | Fetch message history | Backend F05 |
| `/messages/:id/read` | POST | Mark message as read | Backend F06 |
| `/groups` | POST | Create a new group | Backend F04 |
| `/ws` | WS | Real-time messaging | Backend F02/F03/F06/F07 |

*Full API documentation: See `README.md` > API Documentation section.*

## 8. WebSocket Event Schema (Frontend Handling)

| Event Type | Direction | Payload | Action |
|---|---|---|---|
| `send_message` | Out | `{to, content}` | Emit to server |
| `typing_start` | Out | `{target_id}` | Emit when user types |
| `typing_stop` | Out | `{target_id}` | Emit after 3s debounce |
| `message_delivered` | Out | `{message_id}` | Emit when message received |
| `new_message` | In | `Message` | Update Query Cache + Scroll |
| `user_typing` | In | `{username}` | Show "Alice is typing..." header |
| `user_stopped_typing` | In | `{username}` | Hide typing header |
| `receipt_update` | In | `{msg_id, status}` | Update message icon in list |
