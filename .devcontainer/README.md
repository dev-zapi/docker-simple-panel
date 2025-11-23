# Dev Container 配置

本项目包含 VS Code Dev Container 配置，可在 Windows 机器上无需本地 Docker 环境即可进行开发和测试。

## 前置要求

1. **Windows 系统需安装：**
   - [Visual Studio Code](https://code.visualstudio.com/)
   - [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop) 或
   - [WSL 2](https://docs.microsoft.com/en-us/windows/wsl/install) + Docker

2. **VS Code 扩展：**
   - [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

## 使用方法

### 方式 1: 使用命令面板

1. 在 VS Code 中打开项目
2. 按 `F1` 或 `Ctrl+Shift+P` 打开命令面板
3. 输入并选择 `Dev Containers: Reopen in Container`
4. 等待容器构建和启动（首次使用需要几分钟）

### 方式 2: 使用状态栏

1. 在 VS Code 中打开项目
2. 点击左下角的绿色图标（><）
3. 选择 `Reopen in Container`

## 功能特性

Dev Container 已配置以下功能：

- ✅ **Go 1.21 环境**：预装 Go 编译器和工具链
- ✅ **Docker-in-Docker**：容器内可访问 Docker 守护进程
- ✅ **Go 工具**：gopls、delve、staticcheck、golangci-lint
- ✅ **端口转发**：自动转发 8080 端口到本地
- ✅ **代码提示**：完整的 Go 语言服务器支持
- ✅ **格式化**：保存时自动格式化代码
- ✅ **依赖管理**：容器启动时自动下载依赖

## 开发工作流

容器启动后，你可以：

```bash
# 构建项目
go build -o docker-simple-panel .

# 运行应用
./docker-simple-panel

# 运行测试脚本
./test-api.sh

# 使用示例脚本
./example-usage.sh
```

## 环境变量

Dev Container 已预配置以下环境变量：

- `SERVER_PORT=8080`
- `DATABASE_PATH=/workspace/docker-panel.db`
- `JWT_SECRET=dev-secret-key-change-in-production`
- `DOCKER_SOCKET=/var/run/docker.sock`

你可以在 `.devcontainer/docker-compose.yml` 中修改这些变量。

## 安全说明

⚠️ **开发环境安全提示：**

本 Dev Container 配置通过以下方式访问 Docker：
1. **Docker socket 挂载** - 将宿主机的 `/var/run/docker.sock` 挂载到容器中
2. **Docker-in-Docker 特性** - 使用官方的 Docker-in-Docker 功能

**注意事项：**
- 此配置适用于开发和测试环境
- Docker socket 访问赋予容器完全的 Docker 控制权限
- 不要在生产环境中使用此配置
- 建议仅在信任的开发机器上使用
- 非 root 用户（vscode）运行以降低风险

## 故障排除

### 容器无法访问 Docker

确保 Docker Desktop 已启动。Dev Container 配置已经通过 Docker socket (`/var/run/docker.sock`) 挂载方式提供 Docker 访问。

**Windows 用户：**
- 推荐使用 WSL 2 后端（在 Docker Desktop 设置中启用）
- 在 Docker Desktop 设置的 "Resources" -> "WSL Integration" 中启用 WSL 集成
- 确保 WSL 2 已安装并设置为默认版本

**注意：** 本配置使用 Docker socket 挂载和 Docker-in-Docker 特性，已经提供了容器内的 Docker 访问能力，无需额外配置。如果遇到权限问题，确保 Docker Desktop 正在运行且 WSL 集成已启用。

### 端口已被占用

如果 8080 端口已被占用，可以：
1. 修改 `.devcontainer/docker-compose.yml` 中的 `SERVER_PORT` 环境变量
2. 修改 `.devcontainer/devcontainer.json` 中的 `forwardPorts` 配置

### 容器构建失败

尝试重新构建容器：
1. 按 `F1` 打开命令面板
2. 选择 `Dev Containers: Rebuild Container`

## 更多资源

- [VS Code Dev Containers 文档](https://code.visualstudio.com/docs/devcontainers/containers)
- [Dev Container 规范](https://containers.dev/)
