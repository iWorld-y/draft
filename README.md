# Kratos Monorepo Project

这是一个结合了 Go Kratos 后端和 Vite 前端的单体仓库 (Monorepo) 项目。

## 目录结构

- `proto/`: 原始 Protobuf 定义文件 (.proto)，作为项目的接口规范来源。
- `backend/`: Go Kratos 后端项目，包含业务逻辑、数据持久化及生成的 API 代码 (`backend/api/`)。
- `frontend/`: Vite + React + TypeScript 前端项目。
- `docker-compose.yml`: 用于一键启动前后端容器化环境。
- `Makefile`: 根目录任务管理器，封装了常用的开发与构建指令。
- `buf.gen.yaml`: Buf 代码生成配置文件。

## 快速开始

### 1. 初始化环境

安装必要的工具（protoc, buf, kratos, wire, pnpm 等）并下载依赖：

```bash
# 执行根目录初始化，将自动初始化 proto、backend 和 frontend
make init
```

### 2. 开发与运行

```bash
# 同时启动后端 (kratos run) 和前端 (pnpm dev)
make dev

# 仅启动后端
make dev-backend

# 仅启动前端
make dev-frontend
```

### 3. 代码生成 (Protobuf & Wire)

当修改了 `proto/` 下的定义或后端依赖注入逻辑时：

```bash
# 生成 Go API 代码 (pb.go, http, grpc 等)
make api

# 在 backend 目录下手动执行 wire 依赖注入（如果需要）
cd backend && make wire
```

### 4. 构建与部署

```bash
# 编译前后端项目
make build

# 使用 Docker Compose 启动服务
make docker-up

# 停止并移除容器
make docker-down
```

## 后端开发规范 (Kratos)

后端基于 [Kratos](https://go-kratos.dev/) 框架，遵循其标准的分层结构：
- `api/`: 自动生成的 Protobuf 转换代码。
- `internal/biz/`: 业务逻辑层 (Domain/UC)。
- `internal/data/`: 数据访问层 (Repo)。
- `internal/service/`: 接口实现层 (DTO 转换)。
- `internal/server/`: HTTP/gRPC 服务端配置。

更多后端细节请参考 [backend/README.md](backend/README.md)。

