# Vocabulary Learning Monorepo

基于 Go Kratos + React 的单词记忆应用，包含后端服务、前端界面、Proto 接口定义以及 Docker 开发/部署配置。

## 项目结构

- `proto/`: Protobuf 接口定义（单一事实来源）
- `backend/`: Kratos 后端（业务、数据访问、服务实现）
- `frontend/`: Vite + React + TypeScript 前端
- `docker-compose.dev.yml`: 开发环境（热更新）
- `docker-compose.yml`: 生产风格容器编排
- `Makefile`: 根目录常用任务入口

## 环境要求

- Docker / Docker Compose（推荐开发方式）
- Go `1.26+`（本地开发后端时需要）
- Node.js `20+` 与 `pnpm`（本地开发前端时需要）
- Buf（修改 Proto 并生成代码时需要）

## 快速开始

### 1) 初始化依赖

```bash
make init
```

该命令会执行：
- `proto`: `buf dep update`
- `backend`: `buf dep update && go mod download`
- `frontend`: `pnpm install`

### 2) 启动开发环境（推荐）

```bash
make dev
```

`make dev` 使用 `docker-compose.dev.yml` 启动以下服务：
- `postgres`（5432）
- `backend-dev`（8000/9000，`air` 热更新）
- `frontend-dev`（8123，Vite 热更新）

停止开发环境：

```bash
make dev-down
```

### 3) 访问地址

- 前端: `http://localhost:8123`
- 后端 HTTP: `http://localhost:8000`
- 后端 gRPC: `localhost:9000`

## 常用命令

```bash
# 生成 API 代码（proto -> backend/api）
make api

# 仅本地运行后端
make dev-backend

# 仅本地运行前端
make dev-frontend

# 编译前后端
make build

# 启动生产风格容器
make docker-up

# 停止生产风格容器
make docker-down
```

## 本地运行说明

- `make dev-backend` 会执行 `cd backend && kratos run`。
- 当前 `backend/configs/config.yaml` 默认数据库地址是 `postgres` 容器主机名。
- 如果你要在宿主机直接跑后端，需要把数据库连接改为本地可达地址（例如 `localhost`）或保持使用 Docker 网络。

## 代码生成

当修改 `proto/` 后：

```bash
make api
```

当修改后端依赖注入后：

```bash
cd backend && make wire
```

## 更多文档

- 后端说明: `backend/README.md`
- 前端说明: `frontend/README.md`
