.PHONY: help init build run run-dev migrate migrate-create worker test test-api test-db \
        version clean proto proto-all proto-clean proto-help lint vet fmt docker-build \
        docker-run docker-compose-up docker-compose-down install-deps health-check \
        deps generate-docs bench load-test security-check dev

# Конфигурация
APP_NAME = user-service
BIN_DIR = bin
BUILD_INFO = $(shell git describe --tags --always 2>/dev/null || echo "dev")
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROTOC_IMAGE = namely/protoc-all:1.51_2
PROTO_ROOT = pkg/proto
GEN_DIR = pkg/gen

# Главная цель по умолчанию
.DEFAULT_GOAL := help

## 📚 Помощь
help:
	@echo "👤 User Service - Makefile"
	@echo ""
	@echo "Доступные команды:"
	@echo ""
	@echo "📦 Proto файлы:"
	@echo "  make proto              - Build image and generate all proto files"
	@echo "  make proto-generate     - Generate code for internal use"
	@echo "  make proto-pkg          - Generate code for external services"
	@echo "  make proto-pkg-simple   - Simple version for Windows"
	@echo "  make proto-pkg-script   - Generate via script (recommended)"
	@echo "  make proto-clean        - Clean generated files"
	@echo ""
	@echo "🏗️  Сборка и запуск:"
	@echo "  make build              - Сборка бинарника"
	@echo "  make run                - Сборка и запуск сервера"
	@echo "  make run-dev            - Запуск в режиме разработки"
	@echo "  make dev                - Запуск с hot reload (требуется air)"
	@echo "  make clean              - Очистка сборки"
	@echo ""
	@echo "🔧 Управление:"
	@echo "  make migrate            - Выполнить миграции БД"
	@echo "  make migrate-create     - Создать новую миграцию"
	@echo "  make worker             - Запустить фоновых воркеров"
	@echo "  make health-check       - Проверить здоровье сервиса"
	@echo ""
	@echo "🧪 Тестирование:"
	@echo "  make test               - Запуск всех тестов"
	@echo "  make test-api           - Тестирование API"
	@echo "  make test-db            - Тестирование БД"
	@echo "  make bench              - Бенчмарки"
	@echo "  make load-test          - Нагрузочное тестирование"
	@echo "  make lint               - Линтинг кода"
	@echo "  make vet                - Проверка кода"
	@echo "  make fmt                - Форматирование кода"
	@echo "  make security-check     - Проверка безопасности"
	@echo ""

## 📦 Proto файлы
proto: proto-build proto-generate

proto-build:
	@echo "📦 Building protoc-go image..."
	docker build -t $(PROTOC_IMAGE) -f infra/protoc-go.Dockerfile .
	@echo "✅ Docker image built"

proto-generate:
	@echo "🔧 Generating Go code from shared proto files..."
	docker run --rm \
		-v "$(CURDIR):/workspace" \
		$(PROTOC_IMAGE)
	@echo "✅ Proto files generated"

proto-pkg:
	@echo "🚀 Generating for external services (in pkg/gen/)..."
	@mkdir -p pkg/gen
	@echo "Using Docker image: $(PROTOC_IMAGE)"
	@docker run --rm \
		-v "$(CURDIR):/workspace" \
		$(PROTOC_IMAGE) \
		sh -c ' \
			echo "Finding proto files..." && \
			find pkg/proto -name "*.proto" | while read f; do \
				echo "Processing $$f" && \
				protoc -I pkg/proto -I /include \
					--go_out=pkg/gen \
					--go_opt=paths=source_relative \
					--go-grpc_out=pkg/gen \
					--go-grpc_opt=paths=source_relative \
					$$f || exit 1; \
			done && \
			echo "✅ Shared library generated in pkg/gen/" \
		'
	@echo "✅ Shared library generated"

proto-pkg-simple:
	@echo "🚀 Generating for external services (simple version)..."
	@mkdir -p pkg/gen
	@docker run --rm \
		-v "$(CURDIR):/workspace" \
		$(PROTOC_IMAGE) \
		sh -c 'find pkg/proto -name "*.proto" -exec echo "Processing {}" \; -exec protoc -I pkg/proto -I /include --go_out=pkg/gen --go_opt=paths=source_relative --go-grpc_out=pkg/gen --go-grpc_opt=paths=source_relative {} \;'
	@echo "✅ Shared library generated in pkg/gen/"

proto-pkg-script:
	@echo "🚀 Generating via script..."
	@docker run --rm \
		-v "$(CURDIR):/workspace" \
		$(PROTOC_IMAGE) \
		sh -c ' \
			PROTO_ROOT="pkg/proto" && \
			OUTPUT_DIR="pkg/gen" && \
			mkdir -p $$OUTPUT_DIR && \
			find $$PROTO_ROOT -name "*.proto" | while read proto_file; do \
				echo "📁 Processing: $$proto_file" && \
				protoc -I pkg/proto -I /include \
					--go_out=$$OUTPUT_DIR \
					--go_opt=paths=source_relative \
					--go-grpc_out=$$OUTPUT_DIR \
					--go-grpc_opt=paths=source_relative \
					$$proto_file || exit 1; \
			done && \
			echo "✅ Done! Check $$OUTPUT_DIR" \
		'
	@echo "✅ Generated via script"

proto-clean:
	@echo "🧹 Cleaning generated files..."
	@if exist "internal\gen" rmdir /s /q "internal\gen" 2>nul || rm -rf pkg/gen
	@if exist "pkg\gen" rmdir /s /q "pkg\gen" 2>nul || rm -rf pkg/gen
	@echo "✅ Clean complete"

## 🏗️  Сборка и запуск
build:
	@echo "🔨 Building $(APP_NAME)..."
	mkdir -p $(BIN_DIR)
	go build -ldflags="-X 'main.Version=$(BUILD_INFO)' \
		-X 'main.Commit=$(COMMIT_HASH)' \
		-X 'main.BuildDate=$(BUILD_DATE)'" \
		-o $(BIN_DIR)/$(APP_NAME) ./cmd/user-service
	@echo "✅ Build complete: $(BIN_DIR)/$(APP_NAME)"

run: build
	@echo "🚀 Starting User Service server..."
	@echo "Server will be available at: http://localhost:8081"
	@echo "Health check: http://localhost:8081/health"
	@echo ""
	@cd $(BIN_DIR) && ./$(APP_NAME) --config ../config.yaml

run-dev:
	@echo "🚀 Starting in development mode..."
	@echo "For hot reload use: make dev"
	DEBUG=true go run ./cmd/user-service

dev:
	@echo "🔥 Starting with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "⚠ air is not installed. Install: go install github.com/cosmtrek/air@latest"; \
		echo "Running without hot reload..."; \
		make run-dev; \
	fi

## 🔧 Управление
migrate:
	@echo "🔄 Applying migrations..."
	go run ./cmd/migrate

migrate-create:
	@echo "📝 Creating new migration..."
	@read -p "Enter migration name: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	echo "Creating migration: $${timestamp}_$${name}.sql"; \
	echo "-- Migration: $${timestamp}_$${name}" > db/migrations/$${timestamp}_$${name}.sql; \
	echo "✅ Created: db/migrations/$${timestamp}_$${name}.sql"

migrate-create: build
	@echo "📝 Creating migration..."
	@read -p "Enter migration name: " name; \
	echo "Create file: db/migrations/$${name}_up.sql and $${name}_down.sql"

worker: build
	@echo "👷 Starting workers..."
	@cd $(BIN_DIR) && ./$(APP_NAME) worker --workers 3 --queue user_tasks

health-check:
	@echo "❤️  Health checking service..."
	@if curl -s http://localhost:8081/health > /dev/null; then \
		echo "✅ User Service is running"; \
	else \
		echo "❌ User Service is not available"; \
	fi

## 🧪 Тестирование
test: 
	@echo "🧪 Running all tests..."
	go test -v -race ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	@echo "✅ Tests completed"

test-api:
	@echo "🧪 Testing API..."
	@echo "Starting server in background..."
	@go run ./cmd/user-service &
	@SERVER_PID=$$!
	@sleep 3
	@echo "Testing health endpoint..."
	@curl -s http://localhost:8081/health
	@echo ""
	@echo "Testing user endpoint..."
	@curl -s "http://localhost:8081/api/v1/user?id=test"
	@echo ""
	@kill $$SERVER_PID 2>/dev/null || true
	@echo "✅ API tests completed"

test-db:
	@echo "🧪 Testing database..."
	@echo "⚠ Database tests not configured"
	@echo "Configure database connection in config.yaml"

bench:
	@echo "📊 Running benchmarks..."
	go test -bench=. -benchmem ./...

load-test:
	@echo "⚡ Running load tests..."
	@if command -v k6 > /dev/null; then \
		echo "Create scripts/loadtest.js first"; \
	else \
		echo "⚠ k6 is not installed. Install: https://k6.io/docs/getting-started/installation/"; \
	fi

## 🛠️  Code quality
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠ golangci-lint is not installed"; \
	fi

vet:
	@echo "🔎 Checking code with vet..."
	go vet ./...
	@echo "✅ Vet completed"

fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Formatting completed"

security-check:
	@echo "🔒 Security checking..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "⚠ gosec is not installed. Install: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
	fi

## 📋 Утилиты
version: build
	@echo "📋 Version information:"
	@cd $(BIN_DIR) && ./$(APP_NAME) --version 2>/dev/null || echo "Version command not implemented"

generate-docs: build
	@echo "📖 Generating documentation..."
	@echo "⚠ Documentation generation not configured"
	@echo "Implement OpenAPI/Swagger documentation"

install-deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@echo "✅ Dependencies installed"

deps:
	@echo "🔄 Updating dependencies..."
	go get -u ./...
	go mod tidy
	go mod vendor
	@echo "✅ Dependencies updated"

init: install-deps proto
	@echo "✅ Project initialized"

clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BIN_DIR) coverage.out
	rm -rf pkg/gen
	go clean
	@echo "✅ Clean completed"

## 🌐 Dual API (HTTP + gRPC)
run-dual:
	@echo "🚀 Starting in DUAL mode (HTTP:8081 + gRPC:9091)..."
	@echo "HTTP REST: http://localhost:8081"
	@echo "gRPC:      localhost:9091"
	@echo ""
	go run ./cmd/user-service  --grpc-port=9091

test-dual:
	@echo "🧪 Testing DUAL API..."
	@echo "1. Starting server..."
	@make run-dual &
	@SERVER_PID=$$!
	@sleep 3
	@echo ""
	@echo "2. Testing HTTP API..."
	@curl -s http://localhost:8081/health
	@echo ""
	@echo ""
	@echo "3. Testing gRPC client..."
	@echo "⚠ gRPC client not implemented"
	@echo ""
	@echo "4. Testing HTTP Python client..."
	@echo "⚠ Python client not implemented"
	@echo ""
	@echo "✅ Dual API tests completed"
	@kill $$SERVER_PID 2>/dev/null || true

grpc-client:
	@echo "🚀 Running gRPC client..."
	@echo "⚠ gRPC client not implemented"
	@echo "Create scripts/clients/test_grpc_client.go"

http-client:
	@echo "🌐 Running HTTP client..."
	@echo "⚠ HTTP client not implemented"
	@echo "Create scripts/clients/test_http_client.py"

## 🚀 Quick start
quick-start:
	@echo "🚀 Quick start for User Service"
	@echo ""
	@echo "1. Initialize project:"
	@echo "   make init"
	@echo ""
	@echo "2. Generate proto files:"
	@echo "   make proto"
	@echo ""
	@echo "3. Run in development mode:"
	@echo "   make run-dev"
	@echo ""
	@echo "4. Test the service:"
	@echo "   make test-api"
	@echo ""
	@echo "📡 Service will be available at: http://localhost:8081"
