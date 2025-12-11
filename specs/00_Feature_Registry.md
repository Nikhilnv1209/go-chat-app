# BMAD Feature Registry: Chat Backend

**Project**: Go Chat Backend
**Epic**: E01 - Minimum Viable Product

---

## Feature Index

| Feature ID | Name | Story File | Status |
|------------|------|------------|--------|
| **F01** | User Authentication | `stories/1.1_user_auth.story.md` | DONE |
| **F02** | WebSocket Hub | `stories/1.2_websocket_hub.story.md` | DONE |
| **F03** | Direct Messaging | `stories/1.3_direct_messaging.story.md` | DONE |
| **F04** | Group Messaging | `stories/1.4_group_messaging.story.md` | DONE |
| **F05** | Inbox & History | `stories/1.5_inbox_history.story.md` | DONE |
| **F06** | Read Receipts | `stories/1.6_read_receipts.story.md` | DONE |
| **F07** | Typing Indicators | `stories/1.7_typing_indicators.story.md` | DONE |

---

## Story Breakdown

### F01: User Authentication
| Story ID | Description | Test Cases |
|----------|-------------|------------|
| S01.01 | User Registration | T01.01, T01.02 |
| S01.02 | User Login | T01.03, T01.04 |

### F02: WebSocket Hub
| Story ID | Description | Test Cases |
|----------|-------------|------------|
| S02.01 | WebSocket Connection | T02.01, T02.02, T02.03 |
| S02.02 | Presence Management | T02.04 |

### F03: Direct Messaging
| Story ID | Description | Test Cases |
|----------|-------------|------------|
| S03.01 | Send DM to Online User | T03.01 |
| S03.02 | Send DM to Offline User | T03.02 |
| S03.03 | Conversation Sync | T03.03 |

### F04: Group Messaging
| Story ID | Description | Test Cases |
|----------|-------------|------------|
| S04.01 | Create Group | T04.03 |
| S04.02 | Send Group Message | T04.02 |
| S04.03 | Access Control | T04.01 |

### F05: Inbox & History
| Story ID | Description | Test Cases |
|----------|-------------|------------|
| S05.01 | Get Conversation List | T05.01 |
| S05.02 | Get Message History | T05.02, T05.03 |

### F06: Read Receipts
| Story ID | Description | Test Cases |
|----------|-------------|------------|
| S06.01 | Create Receipt on Send | T06.01 |
| S06.02 | Update to DELIVERED | T06.02 |
| S06.03 | Update to READ | T06.03 |

### F07: Typing Indicators
| Story ID | Description | Test Cases |
|----------|-------------|------------|
| S07.01 | Typing Start Event | T07.01 |
| S07.02 | Typing Stop Event | T07.02 |

---

## Requirements Traceability

| Req ID | Description | Feature |
|--------|-------------|---------|
| R01 | Users must register with unique email | F01 |
| R02 | Passwords must be hashed | F01 |
| R03 | JWT tokens expire in 24h | F01 |
| R04 | Real-time message delivery | F02, F03 |
| R05 | Multi-device support | F02 |
| R06 | Message persistence | F03, F04 |
| R07 | Unified MessageReceipt for DM & Group | F03, F04 |
| R08 | Inbox sorted by last_message_at | F05 |
| R09 | Message delivery confirmation (SENT/DELIVERED/READ) | F06 |
| R10 | Receipt status queryable via API | F06 |
| R11 | Real-time typing status broadcast | F07 |
| R12 | Automatic timeout for stale typing indicators | F07 |

---

## Quick Reference Tags
Use these tags in code comments and commit messages:

```
// [F01] User Authentication
// [S01.02] Login Implementation
// [R03] JWT 24h expiry
// [T01.04] Wrong password test
// [F06] Read Receipts
// [F07] Typing Indicators
```
