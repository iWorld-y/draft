GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	# Windows 下的 `find.exe` 与 bash/shell 中的 `find` 不同。
	# 参阅 https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find。
	# 改为使用 git-bash.exe 来运行 find 等命令行工具，因为每个开发者通常都安装了 Git。
	# Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find ../proto -name *.proto)
endif

.PHONY: init
# 初始化环境
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# 生成内部配置 proto 代码
config:
	cd .. && buf generate backend --path backend/internal/conf

.PHONY: api
# 生成 API proto 代码
api:
	cd .. && buf generate proto

.PHONY: build
# 编译项目
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

# 执行 wire 依赖注入
wire:
	cd cmd/backend && wire

.PHONY: generate
# 执行 go generate 并整理 mod 文件
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# 执行所有生成任务
all:
	make api
	make config
	make generate
	make wire

.PHONY: deploy
# 部署项目
deploy:
	docker build -t draft-backend .

.PHONY: run
# 运行项目
run:
	make build
	make deploy
	docker run --rm -p 8124:8000 -p 8125:9000 -v $(PWD)/configs:/data/conf draft-backend

# 显示帮助信息
help:
	@echo ''
	@echo '使用方法:'
	@echo ' make [target]'
	@echo ''
	@echo '指令列表:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
