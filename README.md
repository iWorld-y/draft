# Kratos Monorepo Project

这是一个结合了 Go Kratos 后端和 Vite 前端的单体仓库 (Monorepo) 项目。

## 目录结构

- `proto/`: [新增] 原始 Protobuf 定义文件 (.proto)。
- `api/`: 跨端共享的生成的代码 (Go pb.go, TS 等)。
- `backend/`: Go Kratos 后端项目。
- `frontend/`: Vite/React 前端项目。
- `scripts/`: 自动化脚本。
- `docker-compose.yml`: 一键启动前后端。
- `Makefile`: 根目录任务管理器。

## 快速开始

### 初始化

```bash
# 开启代理以保证 buf 能下载依赖
make init
```

### 开发环境

```bash
# 同时启动前后端
make dev

# 仅启动后端
make dev-backend

# 仅启动前端
make dev-frontend
```

### 生成 API 代码

```bash
make api
```

### Docker 部署

```bash
make docker-up
```
