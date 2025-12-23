# BMAD Feature Specification: MVP Chat Backend

**Epic ID**: 01
**Title**: Minimum Viable Product (Real-time Chat)
**Status**: APPROVED
**Owner**: Product Manager Agent (Antigraivty)

---

## 1. Overview
The goal of this Epic is to deliver a functional "WhatsApp-lite" backend. It bridges the gap between the high-level `PROJECT_SPECS.md` and the actual code. We will implement User Auth, One-on-One Messaging, Group Messaging, and Basic History.

## 2. Architecture Reference
*   **Database**: PostgreSQL (Users, Messages, Groups, MessageReceipts, Conversations)
*   **Transport**: HTTP (Auth/REST) + WebSocket (Events)
*   **Patterns**: Repository Pattern, Hub-and-Spoke for WebSockets.

*(See `specs/03_Technical_Specification.md` for Schema and API details)*

---

## 3. User Stories (Implementation Tasks)

### Story 1.1: User Authentication System [F01]
**As a** new user,
**I want to** register and login with a secure password,
**So that** I can access the chat features.

*   **Context**: We need a `User` table. Passwords must be hashed with bcrypt. Login returns a JWT.
*   **Technical Implementation**:
    *   `POST /auth/register`: Accepts `{username, email, password}`.
    *   `POST /auth/login`: Accepts `{email, password}`. Returns `{token, user}`.
    *   Middleware: `AuthMiddleware` verifies JWT header.
*   **Acceptance Criteria**:
    *   [ ] Valid registration creates a DB row.
    *   [ ] Duplicate email registration fails with 409 Conflict.
    *   [ ] Login with wrong password fails with 401.
    *   [ ] Login returns a valid JWT signed with server secret.

### Story 1.2: WebSocket Hub & Connection [F02]
**As a** logged-in user,
**I want to** establish a persistent WebSocket connection,
**So that** I can receive messages instantly.

*   **Context**: The heart of the real-time system.
*   **Technical Implementation**:
    *   `GET /ws`: Upgrades HTTP to WebSocket. Requires Auth.
    *   **Hub Logic**: Maintain a Map `map[UserID]*Client`.
    *   **Presence**: On connect, set `IsOnline = true` in DB. On disconnect, update `LastSeen`.
*   **Acceptance Criteria**:
    *   [ ] User can connect to `ws://localhost:8080/ws?token=...`.
    *   [ ] Server logs "Client connected".
    *   [ ] Disconnecting updates the `last_seen` timestamp in DB.

### Story 1.3: Direct Messaging (One-on-One) [F03]
**As a** user,
**I want to** send a text message to another user,
**So that** we can communicate privately.

*   **Context**: Uses the Unified `MessageReceipt` model.
*   **Technical Implementation**:
    *   **WS Event**: User sends JSON `{type: "send_message", to: userID, content: "hi"}`.
    *   **Server Action**: Save Message, Save Receipt, Update Conversation, Push to Socket.
*   **Acceptance Criteria**:
    *   [ ] Message saves to Postgres.
    *   [ ] Recipient receives message immediately if online.
    *   [ ] `Conversation` table updates `LastMessageAt`.

### Story 1.4: Group Management & Messaging [F04]
**As a** user,
**I want to** create a group and talk to multiple people,
**So that** we can coordinate together.

*   **Context**: Requires `Group` and `GroupMember` tables.
*   **Technical Implementation**:
    *   `POST /groups`: JSON `{name: "Family", members: [id1, id2]}`.
    *   **WS Event**: `{type: "send_message", group_id: 1, content: "Hello All"}`.
*   **Acceptance Criteria**:
    *   [ ] Group is created in DB.
    *   [ ] Message sent to Group ID 1 is received by all online members.

### Story 1.5: Inbox & History Sync [F05]
**As a** user,
**I want to** see my past conversations and messages,
**So that** I don't lose context when I restart the app.

*   **Context**: REST APIs for fetching data.
*   **Technical Implementation**:
    *   `GET /conversations`: Returns list of chats.
    *   `GET /messages`: Returns paginated history.
*   **Acceptance Criteria**:
    *   [x] Returns correct list of DMs and Groups.
    *   [x] History loads in correct chronological order.

### Story 1.6: Read Receipts [F06]
**As a** user,
**I want to** see if my message has been delivered and read,
**So that** I know the recipient has seen my message.

*   **Context**: Implements the existing `MessageReceipt` model with SENT/DELIVERED/READ statuses.
*   **Technical Implementation**:
    *   Create `MessageReceiptRepository`.
    *   Update `MessageService` to create receipts on send.
    *   Add `POST /messages/:id/read` endpoint.
    *   Handle `message_delivered` WebSocket event.
*   **Acceptance Criteria**:
    *   [ ] Receipt with status SENT is created when message is sent.
    *   [ ] Receipt updates to DELIVERED when receiver acknowledges.
    *   [ ] Receipt updates to READ when receiver marks as read.
    *   [ ] Sender receives `receipt_update` WebSocket event.

### Story 1.7: Typing Indicators [F07]
**As a** user,
**I want to** see when someone is typing a message to me,
**So that** I know they are actively responding.

*   **Context**: Stateless, WebSocket-only feature for real-time UX.
*   **Technical Implementation**:
    *   Handle `typing_start` and `typing_stop` WebSocket events.
    *   Broadcast typing status to conversation participants.
    *   No database persistence required.
*   **Acceptance Criteria**:
    *   [ ] When User A types, User B receives `user_typing` event.
    *   [ ] When User A stops, User B receives `user_stopped_typing` event.
    *   [ ] Works for both DMs and Group chats.

---

## 4. Work Order (for Development Agent)
1.  [x] **Initialize**: Setup Go Modules, Directory Structure, Docker. (Turbo)
2.  [x] **Database**: Run Migrations (GORM AutoMigrate).
3.  [x] **Auth Module**: Implement Story 1.1.
4.  [x] **WS Module**: Implement Story 1.2.
5.  [x] **Messaging Logic**: Implement Story 1.3 and 1.4.
6.  [x] **REST API**: Implement Story 1.5.
7.  [x] **Read Receipts**: Implement Story 1.6.
8.  [x] **Typing Indicators**: Implement Story 1.7.
