.PHONY: help build run test lint clean docker-build docker-run migrate-up migrate-down sqlc-generate swagger-generate

# Variables
APP_NAME=gopilot
DOCKER_IMAGE=gopilot
DOCKER_TAG=latest
DB_DSN=postgres://postgres:postgres@localhost:5432/gopilot?sslmode=disable

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/$(APP_NAME) ./cmd/server

run: ## Run the application
	go run ./cmd/server/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	go mod download
	go mod tidy

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

migrate-up: ## Run database migrations up
	migrate -path db/migrations -database "$(DB_DSN)" up

migrate-down: ## Run database migrations down
	migrate -path db/migrations -database "$(DB_DSN)" down

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	migrate create -ext sql -dir db/migrations -seq $(NAME)

sqlc-generate: ## Generate sqlc code
	sqlc generate

swagger-generate: ## Generate Swagger documentation
	swag init -g cmd/server/main.go -o docs

gen:openapi: swagger-generate ## Alias for swagger-generate

install-tools: ## Install development tools
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

generate: sqlc-generate swagger-generate ## Generate all code (sqlc + swagger)

dev: ## Run development environment with docker-compose
	docker-compose up -d

dev-down: ## Stop development environment
	docker-compose down

db-reset: migrate-down migrate-up ## Reset database

seed: ## Seed database with sample data
	@echo "Database seeding not yet implemented"
	@echo "To add seed data, create a seed script in db/seeds/"

.DEFAULT_GOAL := help
