# Makefile for the cut_utility project

.PHONY: all fmt lint test coverage clean

# Default target runs all checks
all: fmt lint test

# fmt checks if all go files are formatted with gofmt
fmt:
	@echo "Running gofmt check..."
	@test -z "$(shell gofmt -l .)" || (echo "gofmt found files that need formatting, please run 'gofmt -w .'"; exit 1)
	@echo "gofmt check passed."

# lint runs golangci-lint
lint:
	@echo "Running golangci-lint..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run ./...
	@echo "golangci-lint check passed."

# test runs all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# coverage runs tests and shows coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

# clean removes coverage file
clean:
	@echo "Cleaning up..."
	@rm -f coverage.out


