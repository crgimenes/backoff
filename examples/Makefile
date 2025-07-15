# Makefile for backoff examples

.PHONY: all basic without-jitter http-retry reset test-examples clean help

# Default target
all: test-examples

# Run all examples
test-examples: basic without-jitter http-retry reset

# Individual example targets
basic:
	@echo "=== Running Basic Example ==="
	@cd examples/basic && go run main.go
	@echo

without-jitter:
	@echo "=== Running Without Jitter Example ==="
	@cd examples/without-jitter && go run main.go
	@echo

http-retry:
	@echo "=== Running HTTP Retry Example ==="
	@cd examples/http-retry && go run main.go
	@echo

reset:
	@echo "=== Running Reset Example ==="
	@cd examples/reset && go run main.go
	@echo

# Build all examples to check for compilation errors
build-examples:
	@echo "=== Building All Examples ==="
	@for dir in examples/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			echo "Building $$dir..."; \
			cd "$$dir" && go build -o example main.go && rm -f example && cd ../..; \
		fi \
	done
	@echo "All examples built successfully!"

# Clean any build artifacts
clean:
	@find examples -name "example" -type f -delete
	@echo "Cleaned build artifacts"

# Help target
help:
	@echo "Available targets:"
	@echo "  all           - Run all examples (default)"
	@echo "  test-examples - Run all examples"
	@echo "  basic         - Run basic usage example"
	@echo "  without-jitter- Run without jitter example"
	@echo "  http-retry    - Run HTTP retry example"
	@echo "  reset         - Run reset functionality example"
	@echo "  build-examples- Build all examples to check compilation"
	@echo "  clean         - Remove build artifacts"
	@echo "  help          - Show this help message"
