# SchemaFlow Makefile

.PHONY: test test-coverage test-quick test-verbose test-race test-bench clean build install lint help

# Default target
all: test

# Run tests with coverage
test:
	@echo "Running tests with coverage..."
	@go test -cover ./...

# Run tests with detailed coverage report
test-coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out
	@echo "\nTotal coverage:"
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}'

# Generate HTML coverage report
test-coverage-html: test-coverage
	@echo "Generating HTML coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

# Quick test run (short tests only)
test-quick:
	@echo "Running quick tests..."
	@go test -short ./...

# Verbose test output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@go test -race ./...

# Run benchmarks
test-bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Run specific test
test-run:
	@read -p "Enter test name pattern: " name; \
	echo "Running tests matching $$name..."; \
	go test -v -run "$$name" ./...

# Clean test cache
test-clean:
	@echo "Cleaning test cache..."
	@go clean -testcache
	@rm -f coverage.out coverage.html

# Build the library
build:
	@echo "Building schemaflow..."
	@go build ./...

# Install the library
install:
	@echo "Installing schemaflow..."
	@go install ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run all checks (format, vet, lint, test)
check: fmt vet lint test

# Clean all generated files
clean: test-clean
	@echo "Cleaning generated files..."
	@rm -rf coverage.* *.test *.out

# Show dependencies
deps:
	@echo "Showing dependencies..."
	@go list -m all

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# Verify dependencies
deps-verify:
	@echo "Verifying dependencies..."
	@go mod verify

# Download dependencies
deps-download:
	@echo "Downloading dependencies..."
	@go mod download

# Show module information
info:
	@echo "Module: github.com/monstercameron/schemaflow"
	@echo "Go Version: $$(go version)"
	@echo "Current Coverage: $$(go test -cover ./... 2>&1 | grep -o '[0-9.]*%' | head -1)"
	@echo "Files: $$(find . -name '*.go' -not -path './vendor/*' | wc -l | tr -d ' ')"
	@echo "Lines of Code: $$(find . -name '*.go' -not -path './vendor/*' -exec wc -l {} + | tail -1 | awk '{print $$1}')"

# Help target
help:
	@echo "SchemaFlow Makefile Targets:"
	@echo ""
	@echo "Testing:"
	@echo "  make test                - Run tests with coverage"
	@echo "  make test-coverage       - Generate detailed coverage report"
	@echo "  make test-coverage-html  - Generate HTML coverage report"
	@echo "  make test-quick          - Run quick tests only"
	@echo "  make test-verbose        - Run tests with verbose output"
	@echo "  make test-race           - Run tests with race detection"
	@echo "  make test-bench          - Run benchmarks"
	@echo "  make test-run            - Run specific test by name"
	@echo "  make test-clean          - Clean test cache and reports"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt                 - Format code"
	@echo "  make vet                 - Run go vet"
	@echo "  make lint                - Run linter"
	@echo "  make check               - Run all checks"
	@echo ""
	@echo "Build:"
	@echo "  make build               - Build the library"
	@echo "  make install             - Install the library"
	@echo "  make clean               - Clean generated files"
	@echo ""
	@echo "Dependencies:"
	@echo "  make deps                - Show dependencies"
	@echo "  make deps-update         - Update dependencies"
	@echo "  make deps-verify         - Verify dependencies"
	@echo "  make deps-download       - Download dependencies"
	@echo ""
	@echo "Other:"
	@echo "  make info                - Show module information"
	@echo "  make help                - Show this help message"