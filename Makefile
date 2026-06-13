.PHONY: build clean

build:
	@echo "Building tsctl..."
	@mkdir -p bin/
	@go build -o bin/tsctl cmd/server/main.go
	@echo "Build complete: bin/tsctl"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ docs/
	@echo "Clean complete"

