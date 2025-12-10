# Architecture Decision Record (ADR): Chat Backend

**ADR ID**: ADR-001
**Title**: Core Architecture Decisions
**Status**: APPROVED
**Date**: 2024-12-07
**Owner**: System Architect Agent

---

## 1. Context
We are building a real-time chat backend in Go. Key challenges:
*   Handling thousands of concurrent WebSocket connections.
*   Ensuring message delivery consistency.
*   Maintaining clean, testable code.

---

## 2. Decisions

### ADR-001.1: Layered Architecture
**Decision**: Use a 4-layer architecture: Handler → Service → Repository → Database.

**Rationale**:
*   **Handlers**: Parse HTTP/WS requests, validate input, call services.
*   **Services**: Contain business logic (e.g., "Can this user send to this group?").
*   **Repositories**: Abstract database operations. Return Go structs, not SQL rows.
*   **Models**: GORM structs with validation tags.

**Consequence**: Clear separation. Services can be unit-tested by mocking repositories.

---

### ADR-001.2: Interface-First Design
**Decision**: Define interfaces for all repositories and services before implementation.

**Rationale**:
*   Enables dependency injection.
*   Allows mock implementations for testing.
*   Forces thinking about contracts before code.

**Example**:
```go
type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
    UpdateOnlineStatus(userID uint, isOnline bool) error
}
```

---

### ADR-001.3: WebSocket Hub with Mutex Protection
**Decision**: The Hub struct uses `sync.RWMutex` to protect the client map.

**Rationale**:
*   Multiple goroutines read/write the map (connect, disconnect, broadcast).
*   Go maps are not concurrent-safe by default.

**Implementation**:
```go
type Hub struct {
    clients   map[uint][]*Client // userID -> list of connections
    mu        sync.RWMutex
    register  chan *Client
    unregister chan *Client
    broadcast chan *Message
}
```

---

### ADR-001.4: Configuration via Environment Variables
**Decision**: All config loaded from `.env` file using `os.Getenv()` or `godotenv`.

**Rationale**:
*   12-factor app compliance.
*   Easy to change settings between dev/prod without code changes.
*   Secrets (JWT_SECRET, DB_PASSWORD) never hardcoded.

---

### ADR-001.5: Centralized Error Handling
**Decision**: Define custom error types and a standard API error response.

**Rationale**:
*   Consistent error responses across all endpoints.
*   Easier debugging with error codes.

**Standard Error Response**:
```json
{
  "error": {
    "code": "AUTH_INVALID_CREDENTIALS",
    "message": "Email or password is incorrect"
  }
}
```

---

### ADR-001.6: Middleware Chain
**Decision**: Use a middleware chain for cross-cutting concerns.

**Order**:
1.  **Logger**: Log all requests.
2.  **CORS**: Handle cross-origin requests.
3.  **Auth**: Validate JWT, inject `userID` into context.
4.  **Handler**: Actual endpoint logic.

---

### ADR-001.7: Database Connection Pooling
**Decision**: Configure GORM with connection pool settings.

**Rationale**:
*   Prevents "too many connections" errors under load.
*   Reuses connections for efficiency.

**Implementation**:
```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

### ADR-001.8: Read Receipts with Three-State Model [F06]
**Decision**: Implement read receipts using a three-state model: SENT → DELIVERED → READ.

**Rationale**:
*   **SENT**: Created when message is saved to database.
*   **DELIVERED**: Updated when receiver's client acknowledges receipt via WebSocket.
*   **READ**: Updated when user explicitly views the message (REST endpoint).

**Implementation**:
*   `MessageReceipt` table with `status` field.
*   WebSocket event `message_delivered` for client acknowledgment.
*   REST endpoint `POST /messages/:id/read` for marking as read.
*   Broadcast `receipt_update` event to sender when status changes.

**Consequence**: Provides WhatsApp-like delivery confirmation UX.

---

### ADR-001.9: Stateless Typing Indicators [F07]
**Decision**: Implement typing indicators as stateless WebSocket events with no database persistence.

**Rationale**:
*   Typing status is ephemeral and doesn't need historical tracking.
*   Reduces database writes significantly (typing events are frequent).
*   Simpler implementation with no cleanup logic needed.

**Implementation**:
*   Client sends `typing_start` and `typing_stop` events.
*   Server broadcasts to conversation participants immediately.
*   No state stored in Hub or database.
*   Client responsible for auto-stopping after 3 seconds of inactivity.

**Consequence**:
*   **Positive**: Zero database overhead, instant propagation.
*   **Negative**: If server restarts, typing state is lost (acceptable for this feature).

---

## 3. Alternatives Considered

| Decision | Alternative | Why Rejected |
|----------|-------------|--------------|
| Layered Architecture | Hexagonal | Overkill for MVP |
| GORM | Raw SQL | Slower development |
| Mutex in Hub | sync.Map | Less control, harder to debug |
| godotenv | Viper | Simpler for MVP |
| Read Receipts: 3-state | 2-state (Sent/Read) | Less granular, doesn't show delivery |
| Typing Indicators: Stateless | Store in Redis | Unnecessary complexity for MVP |

---

## 4. Consequences
*   **Positive**: Clear structure, testable, maintainable, feature-rich UX.
*   **Negative**: Slight overhead from abstraction layers.
