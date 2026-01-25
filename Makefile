.PHONY: build docs clean run

build:
	@echo "Building tsctl..."
	@mkdir -p bin/
	@go build -o bin/tsctl cmd/server/main.go
	@echo "Build complete: bin/tsctl"

docs:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger docs generated in docs/"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ docs/
	@echo "Clean complete"

run:
	@go run cmd/server/main.go

install-deps:
	@echo "Installing dependencies..."
	@go get -u github.com/gin-gonic/gin
	@go get -u github.com/sirupsen/logrus
	@go get -u gopkg.in/yaml.v3
	@go get -u github.com/swaggo/swag/cmd/swag
	@go get -u github.com/swaggo/gin-swagger
	@go get -u github.com/swaggo/files
	@go get -u github.com/natefinch/lumberjack
	@go mod tidy
	@echo "Dependencies installed"
