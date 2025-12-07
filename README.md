# Go Chat App (Backend)

A high-performance, real-time chat application backend built with Go (Golang).

## 🚀 Features
- **Real-time Messaging**: Powered by WebSockets for instant communication.
- **Direct Messaging**: One-on-one private conversations with unified receipt tracking.
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

### Real-time Communication (WebSocket)

#### Connect to Hub
- **Endpoint**: `GET /ws`
- **Query Parameter**: `token=<YOUR_JWT_TOKEN>`
- **Description**: Upgrades the HTTP connection to a WebSocket connection. Authenticates user via JWT.

**Events (Client -> Server):**
- **Send Message**:
  ```json
  {
    "type": "send_message",
    "payload": {
      "to_user_id": "uuid-string-of-recipient",
      "content": "Hello, World!"
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

## ❓ Troubleshooting

- **Database Connection Failed**:
  - Check if PostgreSQL is running.
  - Verify credentials in `.env` match your database setup.
  - Ensure port `5432` is not blocked.

- **WebSocket Connection Failed**:
  - Ensure you are passing a valid JWT token in the query parameter `?token=...`.
  - Check if the server logs show any authentication errors.

## 📝 Future Roadmap
- [ ] **Group Messaging**: Create groups and broadcast messages.
- [ ] **Message History API**: REST endpoints to fetch conversation history.
- [ ] **Frontend**: React/Vue/Mobile client implementation.
- [ ] **Media Support**: Image and file sharing.

## 📄 License
This project is licensed under the MIT License.
