# PDF Parser

智能PDF文档解析平台 - 将PDF转换为Markdown或纯文本格式

## 功能特点

- 📄 PDF上传与解析
- 🔄 任务状态实时跟踪
- 📝 支持Markdown和纯文本输出格式
- 🎨 现代化Web界面
- ⚡ 异步处理，高效解析

## 技术栈

### 前端
- Vue 3 + TypeScript
- Tailwind CSS
- Axios
- Pinia (状态管理)

### 后端
- Go + Gin
- MySQL + Redis
- MinerU API (PDF解析)

## 快速开始

### 前置要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0
- Redis 7.0

### 后端启动

```bash
cd backend
go mod download
go run ./cmd/server/main.go
```

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

### Docker部署

```bash
docker-compose up -d
```

## API接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/tasks | 创建解析任务 |
| GET | /api/v1/tasks | 获取任务列表 |
| GET | /api/v1/tasks/:id | 获取任务详情 |
| GET | /api/v1/tasks/:id/status | 查询任务状态 |
| GET | /api/v1/tasks/:id/result | 获取解析结果 |
| GET | /api/v1/tasks/:id/download | 下载结果文件 |
| DELETE | /api/v1/tasks/:id | 删除任务 |

## 配置

配置文件位于 `backend/configs/config.yaml`

## 许可证

MIT License
