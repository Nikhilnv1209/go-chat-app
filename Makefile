.PHONY: run dev up down build logs test clean

# Start full stack (database + backend) using podman-compose
run:
	@echo "🚀 Starting full stack with podman-compose..."
	podman-compose up --build

# Start in detached mode
up:
	@echo "🚀 Starting services in background..."
	podman-compose up -d

# Stop all services
down:
	@echo "🛑 Stopping all services..."
	podman-compose down

# Build containers without starting
build:
	@echo "🔨 Building containers..."
	podman-compose build

# View logs
logs:
	podman-compose logs -f

# Local development (requires PostgreSQL running separately)
dev:
	@echo "🔧 Running backend locally (ensure PostgreSQL is running)..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./... -v

# Clean up volumes and containers
clean:
	@echo "🧹 Cleaning up..."
	podman-compose down -v
