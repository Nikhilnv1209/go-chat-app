.PHONY: up down run test

up:
	podman-compose up -d

down:
	podman-compose down

run:
	go run cmd/server/main.go

test:
	go test ./... -v
