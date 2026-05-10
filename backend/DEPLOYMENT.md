# PDF Parser 部署指南

## 快速开始

### 1. 一键构建完整应用

```bash
cd backend
make all
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，填入数据库、Redis 等配置
```

### 3. 运行应用

```bash
./pdf-parser
```

访问 http://localhost:3000 即可使用。

## 构建说明

### 构建方式

```bash
# 构建完整应用（前端+后端+嵌入）
make all

# 或分步构建
make frontend  # 构建前端
make embed    # 复制前端资源到 static 目录
make build    # 构建 Go 后端并嵌入前端
```

### 输出文件

- `pdf-parser` - 单个可执行二进制文件（包含前端资源）
- 大小约 12MB（压缩后更小）

### 部署方式

#### 方式一：直接部署二进制文件

```bash
# 1. 复制二进制文件到服务器
scp pdf-parser user@server:/path/to/app/

# 2. 创建 .env 配置文件
scp .env.example user@server:/path/to/app/.env

# 3. 编辑 .env 配置数据库等信息
ssh user@server
vim /path/to/app/.env

# 4. 运行
cd /path/to/app
./pdf-parser
```

#### 方式二：Docker 部署

```bash
# 构建 Docker 镜像
make docker-build

# 运行容器
make docker-run

# 或手动运行
docker run -p 3000:3000 --env-file .env pdf-parser:latest
```

#### 方式三：使用 systemd 服务

创建 `/etc/systemd/system/pdf-parser.service`:

```ini
[Unit]
Description=PDF Parser Service
After=network.target

[Service]
Type=simple
User=pdfparser
WorkingDirectory=/opt/pdf-parser
ExecStart=/opt/pdf-parser/pdf-parser
Restart=always
EnvironmentFile=/opt/pdf-parser/.env

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable pdf-parser
sudo systemctl start pdf-parser
```

## 环境变量配置

详细配置项请参考 [.env.example](.env.example) 文件。

### 必需配置

```bash
# 数据库 (MySQL)
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_NAME=pdf_parser

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 可选配置

```bash
# MinerU API (PDF 解析服务)
MINERU_API_KEY=your_api_key

# OAuth 认证 (AAA 认证中心)
OAUTH_CLIENT_ID=your_client_id
OAUTH_CLIENT_SECRET=your_client_secret
```

## 目录结构

```
.
├── pdf-parser          # 可执行二进制文件
├── .env               # 环境变量配置
├── configs/           # 配置文件
│   └── config.yaml
├── uploads/           # 上传文件目录 (运行时创建)
└── logs/              # 日志目录 (运行时创建)
```

## 健康检查

```bash
curl http://localhost:3000/health
```

## API 文档

访问 http://localhost:3000/api/v1 查看 API 文档。

## 故障排除

### 数据库连接失败

1. 检查 MySQL 是否运行
2. 验证 .env 中的数据库配置
3. 确保数据库已创建

### Redis 连接失败

1. 检查 Redis 是否运行
2. 验证 .env 中的 Redis 配置

### 前端资源加载失败

1. 重新构建：`make clean && make all`
2. 检查日志输出
