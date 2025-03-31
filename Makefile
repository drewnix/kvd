
.PHONY: build clean test test-cover lint vet fmt run serve all bench install

# Default target
all: lint test build

# Build the application
build:
	go build -o kv -ldflags="-s -w" main.go

# Install the application to GOPATH/bin
install:
	go install -ldflags="-s -w"

# Run code formatting
fmt:
	go fmt ./...

# Run linting
lint:
	go vet ./...
	@if command -v golint >/dev/null; then \
		golint ./...; \
	else \
		echo "golint not installed, skipping"; \
	fi

# Run static code analysis
vet:
	go vet ./...

# Run the application
run:
	go run main.go

# Run the server
serve:
	go run main.go serve

# Run tests
test:
	go test -v -short ./...

# Run tests with coverage report
test-cover:
	go test -cover -v ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Clean build artifacts
clean:
	rm -f kv
	rm -f coverage.out
	go clean

# Security scan (if gosec is installed)
security:
	@if command -v gosec >/dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not installed, skipping security scan"; \
	fi
