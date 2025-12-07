# BMAD Progress Tracker

**Current Status**: 🟢 Phase 1 - F01 User Auth (Story 1.1) Completed
**Last Action**: Implemented User Registration & Login.
**Next Action**: Implement [F02] WebSocket Hub.

---

## 📊 Epic Status: MVP Chat Backend
| Feature | ID | Status | Completion |
|---------|----|--------|------------|
| **User Authentication** | `F01` | ✅ **DONE** | 100% |
| **WebSocket Hub** | `F02` | 🚧 **IN PROGRESS** | 0% |
| **Direct Messaging** | `F03` | 🔴 TODO | 0% |
| **Group Messaging** | `F04` | 🔴 TODO | 0% |
| **Inbox & History** | `F05` | 🔴 TODO | 0% |

---

## 📅 Detailed Work Log

### [F00] Infrastructure & Setup
| Task | Status | Notes |
|------|--------|-------|
| Setup Project & Go Mod | ✅ DONE | `chat-app` module initialized. |
| Setup Docker/Podman | ✅ DONE | Postgres container running. |
| Define GORM Models | ✅ DONE | Models for User, Message, Group, etc. |
| Database Migration | ✅ DONE | AutoMigrate successful. |

---

### [F01] User Authentication
**Story 1.1: Registration & Login** (`stories/1.1_user_auth.story.md`)
| Task | Status | File(s) |
|------|--------|---------|
| Define Core Interfaces | ✅ DONE | `internal/repository/interfaces.go`, `internal/service/interfaces.go` |
| Define Custom Errors | ✅ DONE | `internal/errors/errors.go` |
| Implement Repository | ✅ DONE | `internal/repository/user_repo.go` |
| Implement JWT Logic | ✅ DONE | `pkg/jwt/jwt.go` |
| Implement Service | ✅ DONE | `internal/service/auth_service.go` |
| Implement Handlers | ✅ DONE | `internal/handlers/auth_handler.go` |
| Wire up in `main.go` | ✅ DONE | `cmd/server/main.go` |
| **Verification** | ✅ DONE | Verified via Curl |

---

### [F02] WebSocket Hub
**Story 1.2: Connection & Hub** (`stories/1.2_websocket_hub.story.md`)
| Task | Status | File(s) |
|------|--------|---------|
| Create Client Struct | ⬜ TODO | `internal/websocket/client.go` |
| Create Hub Struct | ⬜ TODO | `internal/websocket/hub.go` |
| Implement Register/Unregister | ⬜ TODO | `internal/websocket/hub.go` |
| Implement WS Handler | ⬜ TODO | `internal/handlers/ws_handler.go` |
| Wire up HTTP Upgrade | ⬜ TODO | `cmd/server/main.go` |
| **Verification** | ⬜ TODO | Test connection with `wscat` |

---

### [F03] Direct Messaging
**Story 1.3: One-on-One Messaging** (`stories/1.3_direct_messaging.story.md`)
| Task | Status | File(s) |
|------|--------|---------|
| Implement Message Repository | ⬜ TODO | `internal/repository/message_repo.go` |
| Implement Conversation Repository | ⬜ TODO | `internal/repository/conversation_repo.go` |
| Implement Message Service (DM) | ⬜ TODO | `internal/service/message_service.go` |
| Implement WS WritePump | ⬜ TODO | `internal/websocket/client.go` |
| Implement WS ReadPump (Event Loop) | ⬜ TODO | `internal/websocket/client.go` |
| **Verification** | ⬜ TODO | Send message between 2 users |

---

### [F04] Group Messaging
**Story 1.4: Group Management** (`stories/1.4_group_messaging.story.md`)
| Task | Status | File(s) |
|------|--------|---------|
| Implement Group Repository | ⬜ TODO | `internal/repository/group_repo.go` |
| Implement Group Service | ⬜ TODO | `internal/service/group_service.go` |
| Implement Group Handlers | ⬜ TODO | `internal/handlers/group_handler.go` |
| Update Message Service for Broadcast | ⬜ TODO | `internal/service/message_service.go` |
| **Verification** | ⬜ TODO | Create group & broadcast message |

---

### [F05] Inbox & History
**Story 1.5: History Sync** (`stories/1.5_inbox_history.story.md`)
| Task | Status | File(s) |
|------|--------|---------|
| Implement GetConvos Handler | ⬜ TODO | `internal/handlers/chat_handler.go` |
| Implement GetMessages Handler | ⬜ TODO | `internal/handlers/chat_handler.go` |
| Wire up Routes | ⬜ TODO | `cmd/server/main.go` |
| **Verification** | ⬜ TODO | Fetch history via Curl |

---
**Legend**:
✅ DONE | 🚧 IN PROGRESS | ⬜ TODO | 🔴 BLOCKED
