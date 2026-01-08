.PHONY: build build-all install clean test run

# Binary name
BINARY=imp
VERSION?=0.1.0

# Build flags
LDFLAGS=-ldflags "-s -w -X github.com/namhoang/imperium/cmd.Version=$(VERSION)"

# Default target
all: build

# Build for current platform
build:
	go build $(LDFLAGS) -o bin/$(BINARY) .

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY)-windows-amd64.exe .

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) .

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test ./...

# Run the application
run:
	go run . $(ARGS)

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run
