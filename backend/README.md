# 单词记忆应用后端

基于 Golang (Kratos) + PostgreSQL 的单词记忆应用后端服务。

## 功能特性

- **词典管理**: 支持上传 TXT 文件批量导入单词，自动调用翻译 API 补全释义
- **SM-2 记忆算法**: 实现科学的间隔重复算法，优化记忆效果
- **学习任务**: 每日生成学习任务队列，智能安排新学和复习
- **进度跟踪**: 实时查看学习进度和单词掌握状态

## 技术栈

- **框架**: Kratos (go-kratos/kratos v2)
- **数据库**: PostgreSQL 14+
- **ORM**: 标准 database/sql + pq
- **翻译 API**: 有道智云 API

## 快速开始

### 1. 环境准备

确保已安装以下软件：
- Go 1.26+
- PostgreSQL 14+
- make (可选)

### 2. 数据库配置

```bash
# 创建数据库
createdb vocabulary

# 执行迁移脚本
psql -d vocabulary -f backend/migrations/001_init_schema.sql
```

### 3. 配置文件

编辑 `backend/configs/config.yaml`：

```yaml
data:
  database:
    driver: postgres
    source: postgres://user:password@localhost:5432/vocabulary?sslmode=disable

translator:
  youdao:
    app_key: "your_youdao_app_key"
    app_secret: "your_youdao_app_secret"
```

### 4. 启动服务

#### 方式一：本地运行

```bash
cd backend
go run internal/cmd/backend/main.go
```

#### 方式二：Docker Compose（推荐）

```bash
# 启动所有服务（包括 PostgreSQL）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

服务将在以下地址启动：
- HTTP: http://localhost:8000
- gRPC: localhost:9000
- PostgreSQL: localhost:5432

## API 接口

### 词典管理

#### 创建词典
```bash
POST /api/v1/dictionaries
Content-Type: application/json

{
  "name": "TOEFL 核心词汇",
  "description": "托福必背 3000 词"
}
```

#### 获取词典列表
```bash
GET /api/v1/dictionaries
```

#### 上传词典文件
```bash
POST /api/v1/dictionaries/upload
Content-Type: multipart/form-data

file: <TXT 文件>
name: "TOEFL 核心词汇"
description: "托福必背 3000 词"
```

#### 查询上传任务状态
```bash
GET /api/v1/dictionaries/upload/status/{task_id}
```

### 学习功能

#### 获取今日学习任务
```bash
GET /api/v1/learning/today-tasks?dict_id=1&limit=20
```

#### 提交学习结果
```bash
POST /api/v1/learning/submit
Content-Type: application/json

{
  "word_id": 1001,
  "quality": 4,
  "time_spent": 8
}
```

**答题质量说明 (quality)**:
- 0: 完全不认识
- 1: 有印象但想不起来
- 2: 想起来了但很费力
- 3: 有些犹豫但想起来了
- 4: 轻松想起来
- 5: 脱口而出

## 项目结构

```
backend/
├── api/                    # API 定义
├── cmd/server/            # 服务入口
├── internal/
│   ├── biz/               # 业务逻辑层
│   │   ├── entity/        # 实体定义
│   │   └── repo/          # 仓库接口
│   ├── data/              # 数据访问层
│   ├── service/           # 服务层
│   └── server/            # HTTP/gRPC 服务器
├── pkg/
│   ├── algorithm/         # SM-2 算法实现
│   └── translator/        # 翻译 API 封装
├── migrations/            # 数据库迁移脚本
└── configs/               # 配置文件
```

## 测试

```bash
# 运行单元测试
go test ./...

# 运行 SM-2 算法测试
go test ./pkg/algorithm/...
```

## SM-2 算法说明

本项目使用 SuperMemo-2 算法来计算复习间隔：

1. **E-Factor (遗忘因子)**: 基于答题质量动态调整，范围 1.3-2.5
2. **间隔天数**: 
   - 第 1 次复习: 1 天后
   - 第 2 次复习: 6 天后
   - 第 n 次复习: I(n) = I(n-1) × EF
3. **答错处理**: 答错后重置间隔为 1 天，重新开始复习周期

## 数据库设计

### 表结构

- **dictionaries**: 词典表
- **words**: 单词表，包含记忆算法字段
- **learn_records**: 学习记录表
- **upload_tasks**: 上传任务表

详见 `migrations/001_init_schema.sql`

## 开发计划

- [x] 基础项目搭建
- [x] 数据库设计
- [x] SM-2 算法实现
- [x] 词典管理功能
- [x] 学习功能核心
- [ ] WebSocket 进度推送
- [ ] 用户认证
- [ ] 学习统计报表

## 使用 Kratos 命令

### 安装 Kratos
```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

### 自动化初始化 (wire)
```bash
# 安装 wire
go get github.com/google/wire/cmd/wire

# 生成 wire
cd cmd/backend
wire
```

### Docker
```bash
# 构建镜像
docker build -t vocabulary-backend .

# 运行容器
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf vocabulary-backend
```

## License

MIT
