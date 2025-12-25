.PHONY: build run test clean proto run-center run-login run-gateway run-game run-all

build:
	@echo "Building Knights Server..."
	@go build -o bin/server cmd/server/main.go

run:
	@echo "Running Knights Server (Gateway)..."
	@go run cmd/server/main.go -config=configs/server.yaml

run-center:
	@echo "Running Center Service..."
	@./bin/server -config=configs/center.yaml

run-login:
	@echo "Running Login Service..."
	@./bin/server -config=configs/login.yaml

run-gateway:
	@echo "Running Gateway Service..."
	@./bin/server -config=configs/server.yaml

run-game:
	@echo "Running Game Service..."
	@./bin/server -config=configs/game.yaml

run-all:
	@echo "Running All Services..."
	@bash scripts/start_all.sh

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf bin/

proto:
	@echo "Generating protobuf code..."
	@protoc --go_out=. --go_opt=paths=source_relative pkg/proto/*.proto

deps:
	@echo "Installing dependencies..."
	@go mod download

fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Running linter..."
	@golangci-lint run

help:
	@echo "Available commands:"
	@echo "  make build        - Build the server"
	@echo "  make run          - Run gateway service"
	@echo "  make run-center   - Run center service"
	@echo "  make run-login    - Run login service"
	@echo "  make run-gateway  - Run gateway service"
	@echo "  make run-game     - Run game service"
	@echo "  make run-all      - Run all services"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make proto        - Generate protobuf code"
	@echo "  make deps         - Install dependencies"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
