# docker-simple-panel-ui

一个简洁的 Docker 容器管理面板前端应用，基于 Svelte 和 TypeScript 开发。

## 功能特性

- 🔐 **用户认证**：基于 JWT 的登录系统
- 📦 **容器管理**：查看 Docker 容器的运行状态和健康状态
- 🎮 **容器控制**：支持启动、停止、重启容器操作
- 💾 **卷管理**：查看、删除和浏览 Docker 卷
- 🎨 **现代界面**：使用 emoji 和渐变色的美观 UI 设计

## 技术栈

- **框架**：Svelte 5
- **语言**：TypeScript
- **构建工具**：Vite
- **路由**：svelte-spa-router
- **状态管理**：Svelte Stores

## 开发

### 安装依赖

```bash
npm install
```

### 配置环境变量

复制 `.env.example` 为 `.env`：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```env
# 使用 Mock API 进行开发（不需要后端）
VITE_USE_MOCK_API=true

# 开发环境反向代理目标（后端服务器地址）
VITE_API_PROXY_TARGET=http://localhost:3000

# 可选：自定义 API URL（默认使用 /api 配合反向代理）
# VITE_API_URL=/api
```

### 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:5173 启动。

#### 反向代理说明

开发环境下，Vite 开发服务器会自动将 `/api` 请求代理到 `VITE_API_PROXY_TARGET` 指定的后端服务器。这样可以避免 CORS 问题，无需配置后端 CORS。

例如：
- 前端请求: `http://localhost:5173/api/containers`
- 实际代理到: `http://localhost:3000/api/containers`

### 构建生产版本

```bash
npm run build
```

构建产物将生成在 `dist` 目录。

### 类型检查

```bash
npm run check
```

## 生产部署

### 使用 Docker 部署

项目提供了 Dockerfile 和 nginx 配置，支持生产环境的反向代理。

#### 构建 Docker 镜像

```bash
docker build -t docker-simple-panel-ui .
```

#### 运行容器

```bash
docker run -d -p 80:80 docker-simple-panel-ui
```

#### 使用 Docker Compose

如果需要同时运行前端和后端，可以使用 docker-compose：

```bash
docker-compose up -d
```

请根据实际情况修改 `docker-compose.yml` 中的后端服务配置。

### Nginx 反向代理配置

生产环境下，nginx 会将前端的 `/api` 请求代理到后端服务器。

默认配置文件 `nginx.conf` 中的代理地址为 `http://localhost:3000`，你可以根据实际部署情况修改：

```nginx
location /api/ {
    proxy_pass http://your-backend-server:3000/api/;
    # ... 其他配置
}
```

或者在运行容器时通过环境变量配置：

```bash
docker run -d -p 80:80 \
  -e BACKEND_URL=http://backend:3000 \
  docker-simple-panel-ui
```

### 静态文件部署

如果使用其他 Web 服务器（如 Apache、Caddy 等），需要：

1. 将 `dist` 目录的内容部署到 Web 服务器
2. 配置反向代理将 `/api` 路径转发到后端服务器
3. 配置 SPA 路由支持（所有路由都返回 index.html）

## 项目结构

```
src/
├── components/      # 可复用组件
│   ├── Header.svelte
│   ├── PageLayout.svelte
│   ├── ContainerList.svelte
│   └── ContentHeader.svelte
├── pages/          # 页面组件
│   ├── Login.svelte
│   ├── Home.svelte
│   ├── Volumes.svelte
│   ├── VolumeExplorer.svelte
│   ├── Settings.svelte
│   ├── ContainerLogs.svelte
│   └── ContainerDetail.svelte
├── stores/         # 状态管理
│   ├── authStore.ts
│   └── themeStore.ts
├── services/       # API 服务
│   ├── api.ts
│   ├── mockApi.ts
│   └── mockData.ts
├── types/          # TypeScript 类型定义
│   └── index.ts
├── App.svelte      # 根组件和路由配置
└── main.ts         # 应用入口
```

## 页面说明

### 登录页 (`/login`)
- 用户输入用户名和密码登录
- 成功后保存 JWT token 到 localStorage

### 首页 (`/`)
- 展示容器列表
- 显示容器状态（运行中、已停止等）配合 emoji
- 显示容器健康状态
- 提供容器控制按钮（启动/停止/重启）
- 支持按 Compose 项目或标签分组
- 顶部导航栏展示 DSP logo 和用户菜单

### 容器详情 (`/container/:id`)
- 查看容器详细信息
- 环境变量、网络配置、端口映射
- 挂载点和重启策略

### 容器日志 (`/logs/:id`)
- 实时查看容器日志
- WebSocket 实时流式传输
- 支持日志过滤和搜索

### 卷管理 (`/volumes`)
- 查看所有 Docker 卷
- 显示卷的容器关联关系
- 删除未使用的卷

### 卷浏览器 (`/volumes/:name/explorer`)
- 浏览卷内文件和目录
- 查看文件内容
- 下载文件

### 系统设置 (`/settings`)
- 配置 Docker socket 路径
- 设置日志级别
- 配置卷浏览器镜像
- 会话超时设置

## API 接口

应用需要后端提供以下 API 接口：

### 认证
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册

### 容器管理
- `GET /api/containers` - 获取容器列表
- `GET /api/containers/:id` - 获取容器详情
- `POST /api/containers/:id/start` - 启动容器
- `POST /api/containers/:id/stop` - 停止容器
- `POST /api/containers/:id/restart` - 重启容器
- `GET /api/containers/:id/logs/stream` - WebSocket 日志流

### 卷管理
- `GET /api/volumes` - 获取卷列表
- `DELETE /api/volumes/:name` - 删除卷
- `GET /api/volumes/:name/explorer/list` - 列出卷内文件
- `GET /api/volumes/:name/explorer/file` - 获取文件内容
- `GET /api/volumes/:name/explorer/download` - 下载文件

### 系统配置
- `GET /api/config` - 获取系统配置
- `PUT /api/config` - 更新系统配置
- `GET /api/docker/health` - Docker 健康检查

## 许可证

MIT
