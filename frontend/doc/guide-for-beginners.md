# 前端项目入门文档

欢迎来到本项目！这份文档是专门为前端初学者准备的，旨在帮助你快速了解项目的结构和每个文件的作用。

## 项目简介

这是一个基于 **React**、**TypeScript** 和 **Vite** 构建的现代前端项目。它采用了目前业界非常流行的技术栈，具有开发速度快、运行效率高的特点。

## 目录结构概览

项目根目录下的文件和文件夹如下：

```text
draft/
├── doc/                # 项目文档目录（你现在就在这里）
├── src/                # 源代码目录，所有的业务逻辑都在这里
│   ├── App.tsx         # 根组件，页面的主要内容通常从这里开始
│   └── main.tsx        # 入口文件，负责将 React 应用挂载到网页上
├── index.html          # HTML 模板，浏览器的入口
├── package.json        # 项目配置文件，定义了依赖和脚本
├── pnpm-lock.yaml      # 依赖锁定文件，确保团队成员安装的包版本一致
├── pnpm-workspace.yaml # PNPM 工作区配置（用于管理多项目）
├── tsconfig.json       # TypeScript 配置文件
├── tsconfig.node.json  # 专门用于 Vite 配置的 TypeScript 配置
└── vite.config.js      # Vite 构建工具的配置文件
```

---

## 文件详细说明

### 1. `src/` 目录 (Source Code)
这是你工作最频繁的地方。

- **[App.tsx](../src/App.tsx)**:
  - **作用**：这是 React 的根组件。你可以把它想象成整个网页的“主容器”。
  - **初学者须知**：你可以在这里编写 HTML 结构的 JSX 代码，并定义页面的基本逻辑。

- **[main.tsx](../src/main.tsx)**:
  - **作用**：这是 JavaScript/TypeScript 的执行入口。
  - **初学者须知**：它的主要任务是找到 `index.html` 中的 `<div id="root"></div>` 节点，并将 `App.tsx` 里的内容“渲染”进去。通常你不需要修改这个文件。

### 2. 根目录配置文件

- **[index.html](../index.html)**:
  - **作用**：这是整个应用的 HTML 骨架。
  - **初学者须知**：你会发现里面有一个 `<script type="module" src="/src/main.tsx"></script>`，这行代码告诉浏览器去加载我们的 JS 代码。

- **[package.json](../package.json)**:
  - **作用**：项目的“身份证”。
  - **初学者须知**：
    - `scripts`: 定义了运行项目的命令（如 `pnpm run dev`）。
    - `dependencies`: 项目运行需要的库（如 `react`）。
    - `devDependencies`: 仅在开发阶段需要的工具（如 `vite`、`typescript`）。

- **[vite.config.js](../vite.config.js)**:
  - **作用**：Vite 的配置文件。
  - **初学者须知**：Vite 是一个极快的开发服务器和打包工具。这个文件定义了如何处理 React 插件、路径别名等高级配置。

- **[tsconfig.json](../tsconfig.json)**:
  - **作用**：TypeScript 的配置。
  - **初学者须知**：它规定了 TypeScript 编译器如何检查你的代码。如果你看到代码里有红色波浪线，通常是这里的规则在起作用。

### 3. 其他文件

- **`pnpm-lock.yaml`**: 自动生成的文件，记录了每个依赖包的精确版本。**不要手动修改它**。
- **`pnpm-workspace.yaml`**: 如果项目变大，包含多个子项目时，这个文件用于管理它们。

---

## 如何开始开发？

1. **安装依赖**：
   在终端运行 `pnpm install`。
2. **启动项目**：
   运行 `pnpm run dev`。
3. **查看效果**：
   打开浏览器，访问终端显示的地址（通常是 `http://localhost:5173`）。

---

## 如何部署到服务器？

部署是将你的代码变成一个可以被全世界访问的网站的过程。对于本项目，通常分为两步：

### 1. 构建项目 (Build)
在部署之前，你需要将源代码“打包”成浏览器能直接高效运行的静态文件。
- **命令**：运行 `pnpm run build`。
- **结果**：项目根目录下会生成一个 `dist/` 文件夹。
- **注意**：`dist` 文件夹包含了压缩后的 HTML、JS 和 CSS。你只需要把这个文件夹的内容上传到服务器即可。

### 2. 选择部署方式

根据你的需求，有几种常见的部署方式：

#### A. 自动化部署平台（推荐初学者）
这类平台非常简单，只需连接你的 GitHub 仓库，每次推送代码后它们会自动帮你构建并发布。
- **Vercel** / **Netlify** / **GitHub Pages**：都是免费且强大的选择。

#### B. 传统 Web 服务器 (Nginx / Apache)
如果你有自己的 VPS 云服务器：
1. 将 `dist/` 文件夹的内容上传到服务器。
2. 配置 Nginx 指向该目录。
3. **重要**：因为这是单页面应用 (SPA)，你需要配置 Nginx 将所有 404 请求重定向到 `index.html`。

#### C. 使用 Docker 部署 (推荐 VPS 用户)
如果你在 VPS 上安装了 Docker，这是最整洁的部署方式。它可以确保你的开发环境和生产环境完全一致。

1. **准备文件**：确保根目录下有 `Dockerfile`、`.dockerignore` 和 `nginx.conf`。
2. **构建镜像**：
   ```bash
   docker build -t my-react-app .
   ```
3. **运行容器**：
   ```bash
   docker run -d -p 8080:80 --name my-app-container my-react-app
   ```
   现在你可以通过 `http://你的服务器IP:8080` 访问你的应用了。

祝你在前端学习的旅程中玩得开心！如果有任何问题，请随时查看相关文档或询问导师。
