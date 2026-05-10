# =============================================================================
# PDF Parser Dockerfile (多阶段构建)
# =============================================================================

# 第一阶段：构建前端
FROM node:20-alpine AS frontend-builder

WORKDIR /app

COPY ../frontend/package*.json ./
RUN npm install

COPY ../frontend/ ./
RUN npm run build

# 第二阶段：构建 Go 后端
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 安装编译依赖
RUN apk add --no-cache gcc musl-dev

# 复制前端构建产物
COPY --from=frontend-builder /app/dist ./router/static

# 复制 Go 源码
COPY ./ ./

# 下载依赖
RUN go mod download

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pdf-parser .

# 最终阶段：单文件部署
FROM alpine:3.19

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 复制二进制文件
COPY --from=backend-builder /app/pdf-parser .

# 复制配置文件
COPY ./configs/ ./configs/
COPY ./.env.example ./.env

# 创建必要目录
RUN mkdir -p ./uploads ./logs

# 设置环境变量
ENV GIN_MODE=release

# 暴露端口
EXPOSE 3000

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# 启动命令
ENTRYPOINT ["./pdf-parser"]
