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

## 故障排除

### 容器无法访问 Docker

确保 Docker Desktop 已启动，并且在设置中启用了 "Expose daemon on tcp://localhost:2375 without TLS"（不推荐生产环境）。

或者，确保 `/var/run/docker.sock` 可以被挂载到容器中。

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
