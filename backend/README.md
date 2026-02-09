# Kratos 项目模板

## 安装 Kratos
```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

## 创建服务
```bash
# 创建一个模板项目
kratos new server

cd server
# 添加一个 proto 模板
kratos proto add api/server/server.proto
# 生成 proto 代码
kratos proto client api/server/server.proto
# 根据 proto 文件生成服务的源代码
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```

## 通过 Makefile 生成其他辅助文件
```bash
# 下载并更新依赖
make init
# 根据 proto 文件生成 API 文件（包括：pb.go, http, grpc, validate, swagger）
make api
# 生成所有文件
make all
```

## 自动化初始化 (wire)
```bash
# 安装 wire
go get github.com/google/wire/cmd/wire

# 生成 wire
cd cmd/server
wire
```

## Docker
```bash
# 构建镜像
docker build -t <your-docker-image-name> .

# 运行容器
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```
