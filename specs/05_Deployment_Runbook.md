# BMAD Deployment & Operations Runbook

**Epic ID**: 01
**Title**: Deployment & Environment Configuration
**Status**: APPROVED
**Owner**: DevOps Agent (Antigravity)

---

## 1. Prerequisites
*   Podman & Podman Compose installed.
*   Go 1.21+ installed.
*   `make` installed (optional, for convenience).

---

## 2. Environment Variables

Create a `.env` file in the project root:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=chat_db

# JWT
JWT_SECRET=your-super-secret-key-change-in-production

# Server
SERVER_PORT=8080
```

---

## 3. Running the Application

### Step 1: Start Database
```bash
podman-compose up -d
```
*This starts PostgreSQL on port 5432.*

### Step 2: Run Migrations
Migrations are handled automatically by GORM `AutoMigrate` on server start.
*No manual step required.*

### Step 3: Start the Server
```bash
go run cmd/server/main.go
```
Or with Make:
```bash
make run
```

---

## 4. Podman Compose Reference

**File: `docker-compose.yml`**
```yaml
version: '3.8'
services:
  postgres:
    image: docker.io/library/postgres:15-alpine
    container_name: chat_postgres
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: chat_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## 5. Makefile Reference

**File: `Makefile`**
```makefile
.PHONY: up down run test

up:
	podman-compose up -d

down:
	podman-compose down

run:
	go run cmd/server/main.go

test:
	go test ./... -v
```

---

## 6. Verification Checklist
After starting, verify the system is healthy:

| Check | Command / Action | Expected |
|-------|------------------|----------|
| DB Running | `podman ps` | `chat_postgres` is Up |
| Server Logs | Terminal output | "Server started on :8080" |
| Health Check | `curl localhost:8080/health` | `{"status":"ok"}` |
| WS Test | `wscat -c ws://localhost:8080/ws?token=...` | Connection upgraded |

---

## 7. Troubleshooting

| Issue | Solution |
|-------|----------|
| "Connection Refused" on DB | Run `docker-compose up -d` |
| "JWT_SECRET not set" | Ensure `.env` file exists and is sourced |
| Port 8080 in use | Change `SERVER_PORT` in `.env` |

---

## 8. Production Readiness / Future Steps

Before labeling the system as "Production Ready" (Phase 2), the following tasks must be completed:

| Component | Task | Description |
|-----------|------|-------------|
| **Database** | **Switch to `golang-migrate`** | Replace GORM `AutoMigrate` with versioned SQL migrations (`up`/`down` scripts) to ensure deterministic schema changes. |
| **Security** | **Secrets Management** | Move secrets from `.env` to a secure vault (e.g., HashiCorp Vault, AWS Secrets Manager). |
| **Observability** | **Structured Logging** | Replace standard `log` with `zap` or `slog` for JSON logging. |
| **Infrastructure** | **CI/CD Pipeline** | Automate build and test steps on commit. |
