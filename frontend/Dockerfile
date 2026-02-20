# 第一阶段：构建阶段
FROM node:20-alpine AS build

# 设置工作目录
WORKDIR /app

# 安装 pnpm
RUN npm install -g pnpm

# 复制依赖文件
COPY package.json pnpm-lock.yaml ./

# 安装依赖
RUN pnpm install --frozen-lockfile

# 复制所有源代码
COPY . .

# 接收构建参数
ARG VITE_API_BASE_URL

# 如果传入了 VITE_API_BASE_URL，则写入 .env.production 文件，覆盖默认配置
RUN if [ -n "$VITE_API_BASE_URL" ]; then echo "VITE_API_BASE_URL=$VITE_API_BASE_URL" > .env.production; fi

# 执行构建命令
RUN pnpm run build

# 第二阶段：运行阶段（使用 Nginx）
FROM nginx:stable-alpine

# 将构建后的文件从 build 阶段复制到 Nginx 的默认静态资源目录
COPY --from=build /app/dist /usr/share/nginx/html

# 复制自定义 Nginx 配置以支持 SPA 路由
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露 80 端口
EXPOSE 80

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]
