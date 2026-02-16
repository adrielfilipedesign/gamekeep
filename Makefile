.PHONY: build build-gui build-cli clean test run run-gui install help

# Variables
BINARY_GUI=gamekeep-gui
BINARY_CLI=gamekeep
GUI_PATH=./cmd/gamekeep-gui/main.go
CLI_PATH=./cmd/main.go
BUILD_DIR=./build

help: ## Show this help
	@echo "GameKeep - Build Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: build-gui build-cli ## Build both GUI and CLI

build-gui: ## Build the GUI application
	@echo "Building $(BINARY_GUI)..."
	@go build -o $(BINARY_GUI) $(GUI_PATH)
	@echo "✓ GUI build complete: ./$(BINARY_GUI)"

build-cli: ## Build the CLI application
	@echo "Building $(BINARY_CLI)..."
	@go build -o $(BINARY_CLI) $(CLI_PATH)
	@echo "✓ CLI build complete: ./$(BINARY_CLI)"

build-all: ## Build for all platforms
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building GUI..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_GUI)-linux-amd64 $(GUI_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_GUI)-darwin-amd64 $(GUI_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_GUI)-darwin-arm64 $(GUI_PATH)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_GUI)-windows-amd64.exe $(GUI_PATH)
	@echo "Building CLI..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_CLI)-linux-amd64 $(CLI_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_CLI)-darwin-amd64 $(CLI_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_CLI)-darwin-arm64 $(CLI_PATH)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_CLI)-windows-amd64.exe $(CLI_PATH)
	@echo "✓ Cross-platform build complete in $(BUILD_DIR)/"

install: build ## Install binaries to $GOPATH/bin
	@echo "Installing $(BINARY_GUI) and $(BINARY_CLI)..."
	@cp $(BINARY_GUI) $(GOPATH)/bin/
	@cp $(BINARY_CLI) $(GOPATH)/bin/
	@echo "✓ Installed to $(GOPATH)/bin/"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_GUI) $(BINARY_CLI)
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "✓ Clean complete"

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

run: run-gui ## Run the GUI application (default)

run-gui: ## Run the GUI application
	@go run $(GUI_PATH)

run-cli: ## Run the CLI application
	@go run $(CLI_PATH) $(ARGS)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies updated"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Example usage targets
example-add: ## Example: Add a game (CLI)
	@go run $(CLI_PATH) add-game --name "Test Game" --path "/tmp/gamekeep-test"

example-checkpoint: ## Example: Create checkpoint (CLI)
	@go run $(CLI_PATH) checkpoint --game test_game --name "Test Checkpoint" --note "Testing"

example-list: ## Example: List checkpoints (CLI)
	@go run $(CLI_PATH) list --game test_game
