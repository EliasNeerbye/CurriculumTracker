.PHONY: all build server wasm clean run setup test help

# Variables
BINARY_NAME=curriculum-tracker
SERVER_CMD=./cmd/server
WASM_CMD=./cmd/wasm
WEB_DIR=./web
BIN_DIR=./bin

# Default target
all: build

# Build everything
build: server wasm

# Build server binary
server:
	@echo "Building server..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(SERVER_CMD)
	@echo "Server built successfully!"

# Build WebAssembly
wasm:
	@echo "Building WebAssembly..."
	@mkdir -p $(WEB_DIR)
	@GOOS=js GOARCH=wasm go build -o $(WEB_DIR)/main.wasm $(WASM_CMD)
	@if [ -f "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(WEB_DIR)/; \
	elif [ -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WEB_DIR)/; \
	else \
		echo "Error: Could not find wasm_exec.js in Go installation"; \
		exit 1; \
	fi
	@echo "WebAssembly built successfully!"

# Run the application
run: build
	@echo "Starting server..."
	@$(BIN_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -f $(WEB_DIR)/main.wasm
	@rm -f $(WEB_DIR)/wasm_exec.js
	@echo "Clean complete!"

# Set up database (requires PostgreSQL to be running)
setup:
	@echo "Setting up database..."
	@createdb curriculum_tracker 2>/dev/null || echo "Database may already exist"
	@psql curriculum_tracker < database/schema.sql
	@echo "Database setup complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Development server with auto-restart (requires 'air' tool)
dev:
	@if command -v air >/dev/null 2>&1; then \
		echo "Starting development server with auto-restart..."; \
		air; \
	else \
		echo "Installing 'air' for development auto-restart..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Build for production (optimized)
prod: clean
	@echo "Building for production..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BIN_DIR)/$(BINARY_NAME) $(SERVER_CMD)
	@GOOS=js GOARCH=wasm go build -ldflags="-w -s" -o $(WEB_DIR)/main.wasm $(WASM_CMD)
	@if [ -f "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(WEB_DIR)/; \
	elif [ -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WEB_DIR)/; \
	else \
		echo "Error: Could not find wasm_exec.js in Go installation"; \
		exit 1; \
	fi
	@echo "Production build complete!"

# Check if required tools are installed
check-deps:
	@echo "Checking dependencies..."
	@command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Aborting." >&2; exit 1; }
	@command -v psql >/dev/null 2>&1 || { echo "PostgreSQL client (psql) is required but not installed. Aborting." >&2; exit 1; }
	@echo "All dependencies are available!"

# Help
help:
	@echo "Available commands:"
	@echo "  make build    - Build both server and WebAssembly"
	@echo "  make server   - Build only the server"
	@echo "  make wasm     - Build only WebAssembly frontend"
	@echo "  make run      - Build and run the application"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make setup    - Set up the database"
	@echo "  make test     - Run tests"
	@echo "  make dev      - Run development server with auto-restart"
	@echo "  make prod     - Build optimized production version"
	@echo "  make check-deps - Check if required tools are installed"
	@echo "  make help     - Show this help message"
