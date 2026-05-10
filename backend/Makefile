.PHONY: all build run dev frontend clean build-all

# 颜色定义
GREEN  := \033[0;32m
YELLOW := \033[0;33m
BLUE   := \033[0;34m
NC     := \033[0m

# 路径
BINARY_NAME := pdf-parser
FRONTEND_DIR := ../frontend
STATIC_DIR := router/static

all: build-all

frontend:
	@echo "$(BLUE)📦 构建前端...$(NC)"
	@cd $(FRONTEND_DIR) && npm install && npm run build
	@echo "$(GREEN)✅ 前端构建完成$(NC)"

embed: frontend
	@echo "$(BLUE)📁 复制前端资源到 static 目录...$(NC)"
	@rm -rf $(STATIC_DIR)
	@mkdir -p $(STATIC_DIR)
	@cp -r $(FRONTEND_DIR)/dist/* $(STATIC_DIR)/
	@echo "$(GREEN)✅ 前端资源已准备$(NC)"

build: embed
	@echo "$(BLUE)🔨 构建 Go 后端...$(NC)"
	@go build -ldflags="-s -w" -o $(BINARY_NAME) .
	@echo "$(GREEN)✅ 构建完成: ./$(BINARY_NAME)$(NC)"

build-all: frontend embed
	@echo "$(BLUE)🔨 构建完整应用...$(NC)"
	@go build -ldflags="-s -w" -o $(BINARY_NAME) .
	@echo ""
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)🎉 构建成功！$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@echo ""
	@echo "二进制文件: ./$(BINARY_NAME)"
	@echo "运行命令: ./$(BINARY_NAME)"
	@echo ""
	@echo "首次运行需要配置 .env 文件："
	@echo "  cp .env.example .env"
	@echo "  # 编辑 .env 填入数据库、Redis 等配置"
	@echo ""

dev:
	@echo "$(YELLOW)⚠️  开发模式需要分别启动前后端$(NC)"
	@echo "前端: cd $(FRONTEND_DIR) && npm run dev"
	@echo "后端: go run main.go"

run: build
	@echo "$(BLUE)🚀 启动应用...$(NC)"
	@./$(BINARY_NAME)

clean:
	@echo "$(BLUE)🧹 清理构建文件...$(NC)"
	@rm -rf $(BINARY_NAME)
	@rm -rf $(STATIC_DIR)
	@rm -rf .output
	@echo "$(GREEN)✅ 清理完成$(NC)"

# Docker 相关
docker-build:
	@echo "$(BLUE)🐳 构建 Docker 镜像...$(NC)"
	@docker build -t pdf-parser:latest .

docker-run:
	@echo "$(BLUE)🐳 运行 Docker 容器...$(NC)"
	@docker run -p 3000:3000 --env-file .env pdf-parser:latest

# 帮助
help:
	@echo "$(BLUE)PDF Parser 构建脚本$(NC)"
	@echo ""
	@echo "$(YELLOW)可用命令:$(NC)"
	@echo "  make all         - 构建完整应用（前端+后端+嵌入）"
	@echo "  make frontend    - 仅构建前端"
	@echo "  make embed       - 复制前端资源到 static 目录"
	@echo "  make build       - 构建 Go 后端并嵌入前端"
	@echo "  make run         - 构建并运行"
	@echo "  make dev         - 开发模式说明"
	@echo "  make clean       - 清理构建文件"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make docker-run   - 运行 Docker 容器"
	@echo ""
