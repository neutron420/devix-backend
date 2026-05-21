.PHONY: run build test clean tidy dev

APP_NAME=devix-backend
BUILD_DIR=./bin

run:
	go run cmd/server/main.go

dev:
	@echo "Starting in development mode..."
	go run cmd/server/main.go

build:
	@echo "Building production binary..."
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) cmd/server/main.go

test:
	go test -v -race -cover ./...

tidy:
	go mod tidy

clean:
	@echo "Cleaning build artifacts..."
	@go clean -cache 2>/dev/null || true
	@rm -rf bin coverage.out coverage.html 2>/dev/null || del /q /f bin coverage.out coverage.html 2>nul || true
