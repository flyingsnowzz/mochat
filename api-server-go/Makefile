.PHONY: build run test lint swag docker clean migrate

# 项目变量
APP_NAME    := mochat-api-server
BUILD_DIR   := bin
MAIN_FILE   := cmd/server/main.go
DOCKER_IMAGE := mochat-api

# 构建
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/server $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/server"

# 运行
run:
	@echo "Running $(APP_NAME)..."
	go run $(MAIN_FILE)

# 测试
test:
	@echo "Running tests..."
	go test -v -cover ./...

# 代码检查
lint:
	@echo "Running linters..."
	golangci-lint run ./...

# 生成 Swagger 文档
swag:
	@echo "Generating Swagger docs..."
	swag init -g cmd/server/main.go -o api/swagger

# Docker 构建
docker:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# 清理
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -rf api/swagger
	go clean
	@echo "Clean complete"

# 数据库迁移
migrate:
	@echo "Running migrations..."
	go run scripts/migrate.go
