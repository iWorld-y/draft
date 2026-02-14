# GitHub Actions 工作流说明

本项目包含以下工作流：

## 1. CI/CD Pipeline (`ci.yml`)

**触发条件**: Push 到 main/master 分支、Pull Request、Tag 推送

**包含任务**:
- ✅ 后端构建 (Go + Kratos)
- ✅ 前端构建 (Vite + React)
- ✅ 单元测试
- ✅ 代码检查 (golangci-lint)
- ✅ Protobuf 规范检查 (buf lint)
- ✅ Docker 镜像构建与推送

**镜像推送地址**: `ghcr.io/<username>/<repo>/backend` 和 `ghcr.io/<username>/<repo>/frontend`

## 2. Deploy (`deploy.yml`)

**触发条件**: 手动触发、Tag 推送

**支持部署方式**:
- SSH 远程部署
- Docker Compose 部署
- Kubernetes 部署

> ⚠️ 需要根据实际情况配置部署步骤和 Secrets

## 需要配置的 Secrets

| Secret | 说明 |
|--------|------|
| `GITHUB_TOKEN` | 自动提供，无需手动配置 |
| `SSH_HOST` | 部署服务器地址 |
| `SSH_USERNAME` | SSH 用户名 |
| `SSH_PRIVATE_KEY` | SSH 私钥 |
| `KUBECONFIG` | Kubernetes 配置 (Base64 编码) |

## 快速开始

1. 确保代码已推送到 GitHub
2. 进入仓库 Settings → Secrets and variables → Actions
3. 添加所需的 Secrets
4. 推送代码触发 CI，或手动触发部署
