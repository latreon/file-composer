.PHONY: build run clean test build-gui run-gui

# Binary names
CLI_BINARY_NAME=file-compressor
GUI_BINARY_NAME=file-compressor-gui

# Build directory
BUILD_DIR=./build

# Main build target (CLI)
build:
	@echo "Building CLI..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(CLI_BINARY_NAME) ./cmd/file-compressor

# Run the CLI application
run:
	@go run ./cmd/file-compressor

# Build GUI target
build-gui:
	@echo "Building GUI..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(GUI_BINARY_NAME) ./cmd/gui

# Run the GUI application
run-gui:
	@go run ./cmd/gui

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Install the CLI application to $GOPATH/bin
install:
	@echo "Installing CLI..."
	@go install ./cmd/file-compressor

# Install the GUI application to $GOPATH/bin
install-gui:
	@echo "Installing GUI..."
	@go install ./cmd/gui

# Build all applications for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	# Windows
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(CLI_BINARY_NAME)_windows_amd64.exe ./cmd/file-compressor
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(GUI_BINARY_NAME)_windows_amd64.exe ./cmd/gui
	@echo "Built for Windows (amd64)"

	# Linux
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(CLI_BINARY_NAME)_linux_amd64 ./cmd/file-compressor
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(GUI_BINARY_NAME)_linux_amd64 ./cmd/gui
	@echo "Built for Linux (amd64)"

	# macOS
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(CLI_BINARY_NAME)_darwin_amd64 ./cmd/file-compressor
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(GUI_BINARY_NAME)_darwin_amd64 ./cmd/gui
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(CLI_BINARY_NAME)_darwin_arm64 ./cmd/file-compressor
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(GUI_BINARY_NAME)_darwin_arm64 ./cmd/gui
	@echo "Built for macOS (amd64 and arm64)"

# Default target
default: build build-gui 