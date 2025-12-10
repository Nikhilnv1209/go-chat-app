# BMAD Technical Specification: MVP Chat Backend

**Epic ID**: 01
**Title**: Architecture & API Definition
**Status**: APPROVED
**Owner**: System Architect Agent (Antigravity)
**References**: `specs/01_MVP_Feature_Spec.md`, `specs/02_Architecture_Decisions.md`

---

## 1. Project Directory Structure

```
chat/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration struct & loader
│   ├── database/
│   │   └── db.go                # Database connection & pooling
│   ├── handlers/
│   │   ├── auth_handler.go      # [F01] Auth endpoints
│   │   ├── chat_handler.go      # [F05] Inbox/History endpoints
│   │   ├── group_handler.go     # [F04] Group endpoints
│   │   └── ws_handler.go        # [F02] WebSocket upgrade
│   ├── middleware/
│   │   ├── auth.go              # JWT validation
│   │   ├── cors.go              # CORS headers
│   │   └── logger.go            # Request logging
│   ├── models/
│   │   ├── user.go
│   │   ├── message.go
│   │   ├── message_receipt.go
│   │   ├── group.go
│   │   ├── group_member.go
│   │   └── conversation.go
│   ├── repository/
│   │   ├── interfaces.go        # All repository interfaces
│   │   ├── user_repo.go
│   │   ├── message_repo.go
│   │   ├── group_repo.go
│   │   └── conversation_repo.go
│   ├── service/
│   │   ├── interfaces.go        # All service interfaces
│   │   ├── auth_service.go
│   │   ├── message_service.go
│   │   └── group_service.go
│   ├── websocket/
│   │   ├── hub.go               # Client registry
│   │   ├── client.go            # Single connection handler
│   │   └── message_handler.go   # WS event router
│   └── errors/
│       └── errors.go            # Custom error types
├── pkg/
│   └── jwt/
│       └── jwt.go               # Token generation/validation
├── specs/                       # Documentation (this folder)
├── .env.example                 # Template for environment variables
├── docker-compose.yml
├── Makefile
└── go.mod
```

---

## 2. Component Relationships

```
┌─────────────────────────────────────────────────────────────────────┐
│                           HTTP Client                               │
└───────────────────────────────┬─────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Middleware Chain                            │
│  ┌─────────┐  ┌─────────┐  ┌─────────────────┐                      │
│  │ Logger  │→ │  CORS   │→ │ Auth (optional) │                      │
│  └─────────┘  └─────────┘  └─────────────────┘                      │
└───────────────────────────────┬─────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                           Handlers                                  │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌─────────────┐ │
│  │ auth_handler │ │ chat_handler │ │group_handler │ │ ws_handler  │ │
│  └──────┬───────┘ └──────┬───────┘ └──────┬───────┘ └──────┬──────┘ │
└─────────┼────────────────┼────────────────┼────────────────┼────────┘
          │                │                │                │
          ▼                ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                           Services                                  │
│  ┌──────────────┐ ┌──────────────────┐ ┌──────────────┐             │
│  │ auth_service │ │ message_service  │ │group_service │             │
│  └──────┬───────┘ └────────┬─────────┘ └──────┬───────┘             │
└─────────┼──────────────────┼──────────────────┼─────────────────────┘
          │                  │                  │
          ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Repositories                                │
│  ┌──────────┐ ┌──────────────┐ ┌──────────────┐ ┌─────────────────┐ │
│  │ user_repo│ │ message_repo │ │  group_repo  │ │conversation_repo│ │
│  └────┬─────┘ └──────┬───────┘ └──────┬───────┘ └────────┬────────┘ │
└───────┼──────────────┼────────────────┼──────────────────┼──────────┘
        │              │                │                  │
        ▼              ▼                ▼                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         PostgreSQL                                  │
└─────────────────────────────────────────────────────────────────────┘

                    ┌──────────────────┐
                    │   WebSocket Hub  │
                    │  (Singleton)     │
                    │                  │
                    │  clients map     │
                    │  mutex lock      │
                    │  broadcast chan  │
                    └────────┬─────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
         [Client 1]    [Client 2]    [Client N]
```

---

## 3. Interface Definitions

### Repository Interfaces
```go
// internal/repository/interfaces.go

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uint) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    UpdateOnlineStatus(userID uint, isOnline bool, lastSeen time.Time) error
}

type MessageRepository interface {
    Create(msg *models.Message) error
    FindByConversation(userID, targetID uint, msgType string, limit, beforeID int) ([]models.Message, error)
}

type MessageReceiptRepository interface {
    Create(receipt *models.MessageReceipt) error
    UpdateStatus(messageID, userID uint, status string) error
    FindUnreadCount(userID uint) (int, error)
}

type GroupRepository interface {
    Create(group *models.Group) error
    FindByID(id uint) (*models.Group, error)
    GetMembers(groupID uint) ([]models.GroupMember, error)
    IsMember(groupID, userID uint) (bool, error)
    AddMember(groupID, userID uint, role string) error
}

type ConversationRepository interface {
    Upsert(conv *models.Conversation) error
    FindByUser(userID uint) ([]models.Conversation, error)
    IncrementUnread(userID uint, convType string, targetID uint) error
    ResetUnread(userID uint, convType string, targetID uint) error
}
```

### Service Interfaces
```go
// internal/service/interfaces.go

type AuthService interface {
    Register(username, email, password string) (*models.User, error)
    Login(email, password string) (token string, user *models.User, error)
    ValidateToken(tokenString string) (userID uint, error)
}

type MessageService interface {
    SendDirectMessage(senderID, receiverID uint, content string) (*models.Message, error)
    SendGroupMessage(senderID, groupID uint, content string) (*models.Message, error)
    GetHistory(userID, targetID uint, convType string, limit, beforeID int) ([]models.Message, error)
    MarkAsRead(userID uint, messageIDs []uint) error
}

type GroupService interface {
    Create(creatorID uint, name string, memberIDs []uint) (*models.Group, error)
    AddMember(adminID, groupID, newMemberID uint) error
    RemoveMember(adminID, groupID, memberID uint) error
}
```

### WebSocket Hub Interface
```go
// internal/websocket/hub.go

type HubInterface interface {
    Register(client *Client)
    Unregister(client *Client)
    SendToUser(userID uint, message []byte)
    SendToUsers(userIDs []uint, message []byte)
    IsOnline(userID uint) bool
}
```

---

## 4. Configuration Management

```go
// internal/config/config.go

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Port         string // default: 8080
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

type DatabaseConfig struct {
    Host            string
    Port            string
    User            string
    Password        string
    DBName          string
    MaxIdleConns    int // default: 10
    MaxOpenConns    int // default: 100
    ConnMaxLifetime time.Duration
}

type JWTConfig struct {
    Secret     string
    Expiration time.Duration // default: 24h
}

func Load() (*Config, error) {
    // Load from .env or environment variables
}
```

---

## 5. Error Handling Strategy

### Custom Error Types
```go
// internal/errors/errors.go

type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Status  int    `json:"-"` // HTTP status code
}

// Predefined Errors
var (
    ErrInvalidCredentials = &AppError{"AUTH_INVALID_CREDENTIALS", "Email or password is incorrect", 401}
    ErrEmailExists        = &AppError{"AUTH_EMAIL_EXISTS", "Email already registered", 409}
    ErrUnauthorized       = &AppError{"AUTH_UNAUTHORIZED", "Authentication required", 401}
    ErrForbidden          = &AppError{"AUTH_FORBIDDEN", "You don't have permission", 403}
    ErrNotFound           = &AppError{"RESOURCE_NOT_FOUND", "Resource not found", 404}
    ErrValidation         = &AppError{"VALIDATION_ERROR", "Invalid input", 400}
)
```

### Standard Error Response
```json
{
  "error": {
    "code": "AUTH_INVALID_CREDENTIALS",
    "message": "Email or password is incorrect"
  }
}
```

---

## 6. Middleware Definitions

### Auth Middleware
```go
func AuthMiddleware(jwtService jwt.Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization") // "Bearer <token>"
        userID, err := jwtService.Validate(token)
        if err != nil {
            c.AbortWithStatusJSON(401, errors.ErrUnauthorized)
            return
        }
        c.Set("userID", userID)
        c.Next()
    }
}
```

### Logger Middleware
```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        log.Printf("[%s] %s %d %v",
            c.Request.Method,
            c.Request.URL.Path,
            c.Writer.Status(),
            time.Since(start),
        )
    }
}
```

---

## 7. Concurrency Safety

### Hub Mutex Protection
```go
type Hub struct {
    clients    map[uint][]*Client
    mu         sync.RWMutex
    register   chan *Client
    unregister chan *Client
}

func (h *Hub) SendToUser(userID uint, message []byte) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    if clients, ok := h.clients[userID]; ok {
        for _, client := range clients {
            client.send <- message
        }
    }
}

func (h *Hub) Register(client *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()

    h.clients[client.UserID] = append(h.clients[client.UserID], client)
}
```

---

## 8. Database Schema (Schema Definitions)

### Tables

**1. users**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_online BOOLEAN DEFAULT FALSE,
    last_seen TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**2. groups**
```sql
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP
);
```

**3. group_members**
```sql
CREATE TABLE group_members (
    group_id INT,
    user_id INT,
    role VARCHAR(20) DEFAULT 'MEMBER',
    joined_at TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);
```

**4. messages**
```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id INT,
    receiver_id INT NULL,
    group_id INT NULL,
    content TEXT,
    msg_type VARCHAR(20) DEFAULT 'TEXT',
    created_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

**5. message_receipts**
```sql
CREATE TABLE message_receipts (
    id SERIAL PRIMARY KEY,
    message_id INT,
    user_id INT,
    status VARCHAR(20) DEFAULT 'SENT',
    updated_at TIMESTAMP
);
```

**6. conversations**
```sql
CREATE TABLE conversations (
    id SERIAL PRIMARY KEY,
    user_id INT,
    type VARCHAR(10),
    target_id INT,
    last_message_at TIMESTAMP,
    unread_count INT,
    UNIQUE(user_id, type, target_id)
);
```

---

## 9. API & Protocol Definitions

### REST API
*   `POST /auth/register`
*   `POST /auth/login`
*   `POST /groups`
*   `GET /conversations`
*   `GET /messages`
*   `POST /messages/:id/read` **[F06]** - Mark message as read
*   `GET /messages/:id/receipts` **[F06]** - Query receipt status

### WebSocket Events

#### Client → Server
*   `send_message` - Send a DM or group message
*   `message_delivered` **[F06]** - Acknowledge message delivery
*   `typing_start` **[F07]** - User started typing
*   `typing_stop` **[F07]** - User stopped typing

#### Server → Client
*   `new_message` - Incoming message notification
*   `receipt_update` **[F06]** - Receipt status changed (SENT/DELIVERED/READ)
*   `user_typing` **[F07]** - Another user is typing
*   `user_stopped_typing` **[F07]** - Another user stopped typing

---

## 10. WebSocket Event Schemas

### Client Events

#### send_message
```json
{
  "type": "send_message",
  "payload": {
    "conversation_type": "DM",  // or "GROUP"
    "target_id": "uuid",
    "content": "Hello!"
  }
}
```

#### message_delivered [F06]
```json
{
  "type": "message_delivered",
  "payload": {
    "message_id": "uuid"
  }
}
```

#### typing_start [F07]
```json
{
  "type": "typing_start",
  "payload": {
    "conversation_type": "DM",  // or "GROUP"
    "target_id": "uuid"
  }
}
```

#### typing_stop [F07]
```json
{
  "type": "typing_stop",
  "payload": {
    "conversation_type": "DM",
    "target_id": "uuid"
  }
}
```

### Server Events

#### new_message
```json
{
  "type": "new_message",
  "payload": {
    "id": "uuid",
    "sender_id": "uuid",
    "content": "Hello!",
    "created_at": "2025-12-11T00:30:00Z"
  }
}
```

#### receipt_update [F06]
```json
{
  "type": "receipt_update",
  "payload": {
    "message_id": "uuid",
    "user_id": "uuid",
    "status": "READ",  // SENT, DELIVERED, READ
    "updated_at": "2025-12-11T00:30:00Z"
  }
}
```

#### user_typing [F07]
```json
{
  "type": "user_typing",
  "payload": {
    "user_id": "uuid",
    "username": "Alice",
    "conversation_type": "DM",
    "target_id": "uuid"
  }
}
```

#### user_stopped_typing [F07]
```json
{
  "type": "user_stopped_typing",
  "payload": {
    "user_id": "uuid",
    "conversation_type": "DM",
    "target_id": "uuid"
  }
}
```
