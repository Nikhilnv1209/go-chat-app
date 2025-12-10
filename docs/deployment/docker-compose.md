# Deployment Guide

This guide describes how to deploy the Chat Application using Docker and Docker Compose.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed
- [Docker Compose](https://docs.docker.com/compose/install/) (often included with Docker Desktop)

## Architecture

The deployment consists of two containers:
1. **`app`**: The Go backend server (listening on port 8080).
2. **`postgres`**: The PostgreSQL database (listening on port 5432).

## Configuration

Environment variables are defined in the `docker-compose.yml` file under the `app` service.

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Port the server listens on | `8080` |
| `DB_HOST` | Database host name (service name) | `postgres` |
| `DB_USER` | Database username | `user` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `chat_db` |
| `JWT_SECRET` | Secret key for signing tokens | `...` |

**Security Note:** For production, **DO NOT** hardcode passwords in `docker-compose.yml`. Use a `.env` file or Docker Secrets.

## Building and Running

### 1. Build and Start
Run the following command in the root directory:

```bash
docker-compose up --build -d
```
- `--build`: Forces a rebuild of the Go application image.
- `-d`: Runs containers in detached mode (background).

### 2. Verify Deployment
Check running containers:
```bash
docker-compose ps
```

View logs:
```bash
docker-compose logs -f app
```

### 3. Stopping
To stop the application:
```bash
docker-compose down
```

## Production Considerations

1. **Secrets**: Move sensitive values (`DB_PASSWORD`, `JWT_SECRET`) to a `.env` file which is excluded from git.
2. **Volumes**: Ensure the `postgres_data` volume is backed up.
3. **HTTPS**: Use a reverse proxy (like Nginx or Traefik) in front of the `app` container to handle SSL termination.
