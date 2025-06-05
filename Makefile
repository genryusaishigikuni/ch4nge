# Makefile for CH4NGE API

.PHONY: help build run test clean docker-build docker-run docker-stop docker-clean dev dev-fresh prod

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the Go application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose (development)"
	@echo "  docker-stop   - Stop Docker Compose"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  dev           - Start development environment"
	@echo "  dev-fresh     - Start fresh development environment (clean DB)"
	@echo "  prod          - Start production environment"

# Build the application
build:
	go build -o bin/main .

# Run the application locally
run:
	go run .

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Docker commands
docker-build:
	docker build -t ch4nge-api .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-clean:
	docker-compose down -v
	docker system prune -af

# Development environment
dev:
	@echo "Starting development environment..."
	docker-compose up --build

# Fresh development environment (clean database)
dev-fresh:
	@echo "Starting fresh development environment..."
	docker-compose down -v
	docker-compose up --build

# Production environment
prod:
	@echo "Starting production environment..."
	docker-compose --profile production up --build -d

# Database commands
db-migrate:
	@echo "Running database migrations..."
	go run . migrate

db-seed:
	@echo "Seeding database..."
	go run . seed

# Reset database (for development)
db-reset:
	@echo "Resetting database..."
	docker-compose down
	docker volume rm $$(docker volume ls -q | grep postgres) 2>/dev/null || true
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	sleep 10

# Logs
logs:
	docker-compose logs -f

logs-app:
	docker-compose logs -f app

logs-db:
	docker-compose logs -f postgres

# Health check
health:
	curl -f http://localhost:8080/health || echo "Service is down"