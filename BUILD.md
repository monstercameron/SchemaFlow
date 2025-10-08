# Build Instructions

This document explains how to regenerate files that are excluded from git (binaries, databases, generated code).

## Prerequisites

- Go 1.24.6 or later
- Protocol Buffers compiler (`protoc`)
- protoc-gen-go and protoc-gen-go-grpc plugins

## Install Build Tools

```bash
# Install protoc (macOS)
brew install protobuf

# Install Go protobuf plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Ensure plugins are in PATH
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Generate Protobuf Files

```bash
# Generate workflow engine protobuf files
cd workflow/proto
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       workflow.proto external.proto
cd ../..
```

## Build Example Binaries

```bash
# Build external services example
cd examples/external-services
go build -o external-services main.go
cd ../..

# Build workflow engine example
cd examples/workflow-engine
go build -o workflow-engine main.go
cd ../..

# Build smarttodo example
cd examples/smarttodo/cmd/smarttodo
go build -o ../../smarttodo main.go
cd ../../../..

# Build advanced example
cd examples/advanced
go build -o main cmd/main.go
cd ../..
```

## Build All Examples

```bash
# Build all examples at once
./scripts/build-all.sh
```

## Initialize Databases

Database files are created automatically when you run the examples:

```bash
# Run workflow engine (creates workflow.db)
cd examples/workflow-engine
./workflow-engine

# Run HR application (creates hr.db)
cd examples/hr-application
go run main.go
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific package tests
go test ./workflow/...
```

## Clean Build Artifacts

```bash
# Remove all generated files and binaries
make clean

# Or manually:
rm -f coverage.out
rm -f examples/*/main examples/*/*-app* examples/*/smarttodo
rm -f examples/*/data/*.db examples/*/*.db
rm -f workflow/proto/*.pb.go
```

## Continuous Integration

The CI pipeline automatically:
1. Installs dependencies
2. Generates protobuf files
3. Builds all binaries
4. Runs tests
5. Checks code coverage

## Notes

- **Never commit binaries** - they are gitignored for a reason
- **Never commit database files** - they contain test/local data
- **Generated protobuf files** - regenerate from .proto sources
- **Coverage files** - local testing artifacts only

## Troubleshooting

### "protoc: command not found"
Install Protocol Buffers compiler (see Prerequisites)

### "protoc-gen-go: program not found"
Run: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

### Binary size too large
Use build flags to reduce size:
```bash
go build -ldflags="-s -w" -o binary main.go
```

### Database locked errors
Stop any running instances before regenerating databases
