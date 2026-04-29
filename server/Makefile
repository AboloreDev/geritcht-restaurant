include .env
export

.PHONY: help build run dev lint migrate-up migrate-down docker-up docker-down docs-generate

help: 
	@echo "Available commands:"
	@echo "  build       Build the server binary"
	@echo "  run         Run the server"
	@echo "  dev         Run the server in development mode"
	@echo "  lint        Run linters"
	@echo "  format      Format the code"
	@echo "  migrate-up  Apply database migrations"
	@echo "  migrate-down Rollback database migrations"
	@echo "  docker-up   Start running the container"
	@echo "  docker-down   Rollback running of the container"

build: 
	@echo "Building all binaries...."
	@mkdir -p bin
	@for cmd in cmd/*/; do \
    		if [ -d "$$cmd" ]; then \
    			binary=$$(basename $$cmd); \
    			echo "Building $$binary..."; \
    			go build -o bin/$$binary ./$$cmd; \
    		fi \
    done

run: 
	go run ./cmd/api

dev:
	go run ./cmd/api

format:
	@gofmt -s -w .
	@goimports -w .

lint: 
	golangci-lint run ./...

migrate-up-neon:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

migrate-up-local:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/geritcht_db?sslmode=disable" up

migrate-down-neon:
	migrate -path db/migrations -database "$(DATABASE_URL)" down

migrate-down-local:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/geritcht_db?sslmode=disable" down
	
docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down

docs-generate:
	mkdir -p docs
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude .git,docs,docker,db