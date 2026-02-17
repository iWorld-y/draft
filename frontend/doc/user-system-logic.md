# 用户系统逻辑梳理（已升级为真实鉴权）

## 1. 认证架构

- 后端：`bcrypt` 密码哈希 + `JWT Access Token` + `Refresh Token` 轮换。
- 过期时间：Access/Refresh 均为 7 天。
- Refresh Token 存储：`HttpOnly Cookie`（路径 `/api/v1/auth`）。
- 前端本地仅保存 Access Token 与当前用户信息。

## 2. 后端接口

认证接口：

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me`

业务接口（词典/学习）统一要求 `Authorization: Bearer <access_token>`。

## 3. 鉴权与用户注入

- 路由层通过 `withAuth` 解析并校验 Access Token。
- 鉴权成功后将 `user_id` 注入 context。
- Service 层从 context 读取 `user_id`，不再使用硬编码用户 ID。

## 4. 用户数据入库

新增数据库对象（见 `backend/docs/migrations/002_auth_schema.sql`）：

- `users`
- `auth_refresh_tokens`

其中 refresh token 仅存哈希值，不存明文。

## 5. 多用户隔离

- 词典接口按 `user_id` 创建/查询。
- 上传任务状态校验词典归属。
- 学习接口增加词典和单词归属校验，阻止跨用户访问。

## 6. 前端会话流程

- 登录/注册成功后：保存 Access Token 与用户信息。
- 业务请求 401：自动调用 `/auth/refresh`（依赖 Cookie）刷新令牌并重试一次。
- 退出登录：调用 `/auth/logout` 并清理本地会话。
