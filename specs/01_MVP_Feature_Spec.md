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

## 3. User Stories (Implementation Tasks)
All backend stories are detailed in the `specs/stories/backend/` directory.

| Story ID | Feature Name | Story File | Status |
|---|---|---|---|
| **1.1** | User Authentication | [`1.1_user_auth.story.md`](stories/backend/1.1_user_auth.story.md) | DONE |
| **1.2** | WebSocket Hub | [`1.2_websocket_hub.story.md`](stories/backend/1.2_websocket_hub.story.md) | DONE |
| **1.3** | Direct Messaging | [`1.3_direct_messaging.story.md`](stories/backend/1.3_direct_messaging.story.md) | DONE |
| **1.4** | Group Messaging | [`1.4_group_messaging.story.md`](stories/backend/1.4_group_messaging.story.md) | IN_PROGRESS |
| **1.5** | Inbox & History | [`1.5_inbox_history.story.md`](stories/backend/1.5_inbox_history.story.md) | DONE |
| **1.6** | Read Receipts | [`1.6_read_receipts.story.md`](stories/backend/1.6_read_receipts.story.md) | DONE |
| **1.7** | Typing Indicators | [`1.7_typing_indicators.story.md`](stories/backend/1.7_typing_indicators.story.md) | DONE |
| **1.8** | Conversation Management | [`1.8_conversation_management.story.md`](stories/backend/1.8_conversation_management.story.md) | TODO |

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
