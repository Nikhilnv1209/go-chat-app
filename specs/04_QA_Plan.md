# BMAD Quality Assurance (QA) Strategy

**Epic ID**: 01
**Title**: Testing Strategy & Validation Plan
**Status**: APPROVED
**Owner**: QA Lead Agent (Antigravity)
**References**: `specs/01_MVP_Feature_Spec.md`, `specs/03_Technical_Specification.md`

---

## 1. Testing Pyramid Strategy
We will strictly follow the testing pyramid to ensure reliability without slowing down development.

### Level 1: Unit Tests (Go `testing` package)
*   **Target**: 70%+ Coverage on Logic Layers.
*   **Scope**:
    *   `internal/utils`: Password hashing, JWT generation validation.
    *   `internal/models`: Validation hooks (e.g., stopping empty messages).
    *   *Excluded*: We will mock the database layer for unit tests to keep them fast (<10ms).

### Level 2: Integration Tests (Testcontainers / Docker)
*   **Target**: Critical User Flows.
*   **Scope**:
    *   Spin up a **Real Postgres** container.
    *   Test the `Repository` layer: Does `CreateUser` actually save to DB?
    *   Test `Group` logic: Does adding a member actually update the junction table?

### Level 3: End-to-End (E2E) / API Tests
*   **Tool**: Go HttpTest or a script using `wscat`.
*   **Scope**:
    *   **Flow 1**: User A registers -> Logs in -> Connects WS -> Sends msg -> User B receives.
    *   **Flow 2**: Group Creation -> Message Broadcast -> Verify all 3 members got it.

---

## 2. Test Cases (Matrix)

### A. Authentication [F01]
| ID | Test Case | Input | Expected Output |
| :--- | :--- | :--- | :--- |
| T01.01 | Register Strict | Empty email | 400 Bad Request |
| T01.02 | Register Duplicate | Existing email | 409 Conflict |
| T01.03 | Login Success | Correct creds | 200 OK + JWT |
| T01.04 | Login Failure | Wrong password | 401 Unauthorized |

### B. WebSocket Connectivity [F02]
| ID | Test Case | Input | Expected Output |
| :--- | :--- | :--- | :--- |
| T02.01 | Connect No Token | WS Handshake | 401 Unauthorized |
| T02.02 | Connect Valid | Valid Token | Connection Upgraded (101) |
| T02.03 | Ping/Pong | Send Ping | Receive Pong |
| T02.04 | Presence Update | Disconnect WS | `last_seen` timestamp updated in DB |

### C. Direct Messaging [F03]
| ID | Test Case | Scenario | Validation |
| :--- | :--- | :--- | :--- |
| T03.01 | DM Delivery | A sends to B (Online) | B's socket triggers `new_message` event. |
| T03.02 | DM Offline | A sends to B (Offline) | Message saved in DB. No WS error. |
| T03.03 | Conversation Sync | A sends to B | Both A and B have `conversation` row updated. |

### D. Group Logic [F04]
| ID | Test Case | Scenario | Validation |
| :--- | :--- | :--- | :--- |
| T04.01 | Non-Member Access | User C sends to Group AB | 403 Forbidden. |
| T04.02 | Broadcast | User A sends to Group(A,B,C) | B and C receive event. |
| T04.03 | Create Group | POST /groups | Group created, creator has ADMIN role. |

### E. Inbox & History [F05]
| ID | Test Case | Scenario | Validation |
| :--- | :--- | :--- | :--- |
| T05.01 | Get Conversations | GET /conversations | Returns sorted list of chats. |
| T05.02 | Get Messages | GET /messages?target_id=2 | Returns paginated history. |
| T05.03 | Pagination | GET /messages?before_id=100 | Returns only messages with ID < 100. |

---

## 3. Manual Verification Steps (Post-Implementation)
Since we don't have a frontend, we will use **Postman** and **Websocket Client** (e.g., `wscat`).

**Prerequisites**:
1.  Run `make up` (DB).
2.  Run `make run` (Server).

**Walkthrough**:
1.  **POST /register** (Alice) -> *Get ID 1*
2.  **POST /register** (Bob) -> *Get ID 2*
3.  **POST /login** (Alice) -> *Copy Token A*
4.  **POST /login** (Bob) -> *Copy Token B*
5.  **Connect WS** (Alice): `wscat -c ws://localhost:8080/ws?token=TokenA`
6.  **Connect WS** (Bob): `wscat -c ws://localhost:8080/ws?token=TokenB`
7.  **Send (Alice)**: `{"type":"send", "to": 2, "content":"Hi Bob"}`
8.  **Verify**: Bob's terminal shows the JSON message.

---

## 4. CI/CD Requirements
*   `go vet` must pass.
*   `go fmt` must not change files.
*   `go test ./...` must pass.
