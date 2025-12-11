# Go Chat App (Backend)

A high-performance, real-time chat application backend built with Go (Golang).

## 🚀 Features
- **Real-time Messaging**: Powered by WebSockets for instant communication.
- **Direct Messaging**: One-on-one private conversations with unified receipt tracking.
- **Group Messaging**: Create groups and broadcast messages to multiple members.
- **Inbox & History**: REST APIs to fetch conversation list and message history.
- **Read Receipts**: Track message delivery and read status for all conversations.
- **Typing Indicators**: Real-time typing status for enhanced user experience.
- **Authentication**: Secure JWT-based signup and login flow.
- **Data Persistence**: Robust PostgreSQL database integration using GORM.
- **Scalable Architecture**: Refactored to use UUIDs for all primary and foreign keys.
- **Online Status**: Real-time user online/offline status tracking.

## 🛠️ Tech Stack
- **Language**: Go 1.24+
- **Framework**: Gin Web Framework (High performance HTTP web framework)
- **Database**: PostgreSQL (Advanced open source relational database)
- **ORM**: GORM (The fantastic ORM library for Golang)
- **Real-time**: Gorilla WebSocket (A WebSocket implementation for Go)
- **Auth**: JWT (JSON Web Tokens)

## 📋 Prerequisites
Before you begin, ensure you have met the following requirements:
- **Go**: Version 1.24 or higher installed. [Download Go](https://go.dev/dl/)
- **PostgreSQL**: Local installation or via Docker.
- **Docker & Docker Compose** (Optional): For easier database management.
- **Git**: For version control.

## 📂 Project Structure
```
go-chat-app/
├── cmd/
│   └── server/         # Main entry point for the server
├── internal/
│   ├── config/         # Configuration loading
│   ├── database/       # Database connection and migration
│   ├── handlers/       # HTTP and WebSocket handlers
│   ├── models/         # Database models (User, Message, etc.)
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic layer
│   └── websocket/      # WebSocket hub and client logic
├── pkg/
│   └── jwt/            # JWT helper package
├── specs/              # Project specifications and stories
├── .env                # Environment variables (not committed)
└── README.md           # Project documentation
```

## ⚙️ Setup & Installation

> **Looking to deploy?**
> - [Docker Deployment Guide](docs/deployment/docker-compose.md) (Local/VPS)
> - [Free Hosting Guide](docs/deployment/free-tier.md) (Render + Neon)

### 1. Clone the repository
```bash
git clone https://github.com/Nikhilnv1209/go-chat-app.git
cd go-chat-app
```

### 2. Configure Environment Variables
Create a `.env` file in the root directory. You can copy the structure below:

```env
# Server Configuration
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=chat_db
DB_PORT=5432

# Security
JWT_SECRET=your_super_secret_key_change_this_in_production
JWT_EXPIRATION_HOURS=24
```

### 3. Start Database
**Option A: Using Docker (Recommended)**
Run a PostgreSQL container:
```bash
docker run --name chat_postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=yourpassword -e POSTGRES_DB=chat_db -p 5432:5432 -d postgres
```

**Option B: Local PostgreSQL**
Ensure your local PostgreSQL server is running and a database named `chat_db` exists.

### 4. Run the Server
The server will automatically migrate the database schema on startup to ensure tables exist.

```bash
go run cmd/server/main.go
```
The server will start on `http://localhost:8080`.

## 🧪 Running Tests
To run the automated unit tests for services and handlers:

```bash
go test ./...
```

For comprehensive manual testing scenarios (including WebSocket verify), see the [Manual Testing Guide](docs/testing/manual-guide.md).

## 🔌 API Documentation

### Authentication

#### Register a new user
- **Endpoint**: `POST /auth/register`
- **Body**:
  ```json
  {
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
  }
  ```
- **Response**: `201 Created` with User object (excluding password).

#### Login
- **Endpoint**: `POST /auth/login`
- **Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "securepassword123"
  }
  ```
- **Response**: `200 OK`
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { ... }
  }
  ```

### Group Messaging

#### Create a new group
- **Endpoint**: `POST /groups`
- **Headers**: `Authorization: Bearer <YOUR_JWT_TOKEN>`
- **Body**:
  ```json
  {
    "name": "Family",
    "member_ids": ["uuid-1", "uuid-2", "uuid-3"]
  }
  ```
- **Response**: `201 Created`
  ```json
  {
    "id": "group-uuid",
    "name": "Family"
  }
  ```

#### Add member to group
- **Endpoint**: `POST /groups/:id/members`
- **Headers**: `Authorization: Bearer <YOUR_JWT_TOKEN>`
- **Body**:
  ```json
  {
    "user_id": "uuid-of-new-member"
  }
  ```
- **Response**: `200 OK`
  ```json
  {
    "message": "member added successfully"
  }
  ```

### Inbox & History

#### Get Conversations (Inbox)
- **Endpoint**: `GET /conversations`
- **Headers**: `Authorization: Bearer <YOUR_JWT_TOKEN>`
- **Description**: Returns all conversations (DMs and Groups) for the authenticated user, sorted by last message time.
- **Response**: `200 OK`
  ```json
  [
    {
      "id": "conversation-uuid",
      "type": "DM",
      "target_id": "user-uuid",
      "target_name": "Bob",
      "last_message_at": "2023-10-27T10:00:00Z",
      "unread_count": 3
    },
    {
      "id": "conversation-uuid-2",
      "type": "GROUP",
      "target_id": "group-uuid",
      "target_name": "Family",
      "last_message_at": "2023-10-26T09:00:00Z",
      "unread_count": 0
    }
  ]
  ```

#### Get Message History
- **Endpoint**: `GET /messages`
- **Headers**: `Authorization: Bearer <YOUR_JWT_TOKEN>`
- **Query Parameters**:
  - `target_id` (required): UUID of the user (for DM) or group (for GROUP)
  - `type` (optional): `DM` or `GROUP` (defaults to `DM`)
  - `limit` (optional): Maximum number of messages to return (defaults to 50)
- **Description**: Returns message history for a specific conversation. Automatically resets unread count.
- **Example**: `GET /messages?target_id=uuid-of-user&type=DM&limit=20`
- **Response**: `200 OK`
  ```json
  [
    {
      "id": "msg-uuid-1",
      "sender_id": "user-uuid",
      "receiver_id": "your-uuid",
      "content": "Hello!",
      "created_at": "2023-10-27T10:01:00Z"
    },
    {
      "id": "msg-uuid-2",
      "sender_id": "your-uuid",
      "receiver_id": "user-uuid",
      "content": "Hey there!",
      "created_at": "2023-10-27T10:00:00Z"
    }
  ]
  ```

### Real-time Communication (WebSocket)

#### Connect to Hub
- **Endpoint**: `GET /ws`
- **Query Parameter**: `token=<YOUR_JWT_TOKEN>`
- **Description**: Upgrades the HTTP connection to a WebSocket connection. Authenticates user via JWT.

**Events (Client -> Server):**
- **Send Direct Message**:
  ```json
  {
    "type": "send_message",
    "payload": {
      "to_user_id": "uuid-string-of-recipient",
      "content": "Hello, World!"
    }
  }
  ```

- **Send Group Message**:
  ```json
  {
    "type": "send_message",
    "payload": {
      "group_id": "uuid-string-of-group",
      "content": "Hello Team!"
    }
  }
  ```

- **Typing Start** (DM or Group):
  ```json
  {
    "type": "typing_start",
    "payload": {
      "conversation_type": "DM",
      "target_id": "uuid-of-recipient-or-group"
    }
  }
  ```

- **Typing Stop** (DM or Group):
  ```json
  {
    "type": "typing_stop",
    "payload": {
      "conversation_type": "DM",
      "target_id": "uuid-of-recipient-or-group"
    }
  }
  ```

- **Message Delivered** (Acknowledge receipt):
  ```json
  {
    "type": "message_delivered",
    "payload": {
      "message_id": "uuid-of-message"
    }
  }
  ```

**Events (Server -> Client):**
- **New Message**:
  ```json
  {
    "type": "new_message",
    "payload": {
      "id": "msg-uuid",
      "content": "Hello, World!",
      "sender_id": "sender-uuid",
      "created_at": "timestamp"
    }
  }
  ```

- **Message Sent** (Acknowledgment):
  ```json
  {
    "type": "message_sent",
    "payload": {
      "id": "msg-uuid",
      "content": "Hello, World!",
      "sender_id": "your-uuid",
      "created_at": "timestamp"
    }
  }
  ```

- **User Typing**:
  ```json
  {
    "type": "user_typing",
    "payload": {
      "user_id": "uuid-of-typing-user",
      "username": "Alice",
      "conversation_type": "DM",
      "target_id": "conversation-id"
    }
  }
  ```

- **User Stopped Typing**:
  ```json
  {
    "type": "user_stopped_typing",
    "payload": {
      "user_id": "uuid-of-user",
      "conversation_type": "DM",
      "target_id": "conversation-id"
    }
  }
  ```

- **Receipt Update**:
  ```json
  {
    "type": "receipt_update",
    "payload": {
      "message_id": "msg-uuid",
      "user_id": "uuid-who-read-it",
      "status": "READ",
      "updated_at": "timestamp"
    }
  }
  ```

### Read Receipts

#### Mark Message as Read
- **Endpoint**: `POST /messages/:id/read`
- **Headers**: `Authorization: Bearer <YOUR_JWT_TOKEN>`
- **Description**: Marks a message as read by the authenticated user.
- **Response**: `200 OK`
  ```json
  {
    "message": "message marked as read"
  }
  ```

#### Get Message Receipts
- **Endpoint**: `GET /messages/:id/receipts`
- **Headers**: `Authorization: Bearer <YOUR_JWT_TOKEN>`
- **Description**: Retrieves read receipt status for a specific message.
- **Response**: `200 OK`
  ```json
  [
    {
      "id": "receipt-uuid",
      "message_id": "msg-uuid",
      "user_id": "reader-uuid",
      "status": "READ",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
  ```

## ❓ Troubleshooting

- **Database Connection Failed**:
  - Check if PostgreSQL is running.
  - Verify credentials in `.env` match your database setup.
  - Ensure port `5432` is not blocked.

- **WebSocket Connection Failed**:
  - Ensure you are passing a valid JWT token in the query parameter `?token=...`.
  - Check if the server logs show any authentication errors.

## 📝 Future Roadmap
- [x] **Group Messaging**: Create groups and broadcast messages.
- [x] **Inbox & History API**: REST endpoints to fetch conversation list and message history.
- [x] **Read Receipts**: Track message delivery and read status.
- [x] **Typing Indicators**: Real-time typing status for DM and group conversations.
- [ ] **Frontend**: React/Vue/Mobile client implementation.
- [ ] **Media Support**: Image and file sharing.
- [ ] **Voice/Video Calling**: WebRTC integration for call support.

## 📄 License
This project is licensed under the MIT License.
