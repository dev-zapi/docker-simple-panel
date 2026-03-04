<div align="center">

# Docker Simple Panel

**A modern, lightweight dashboard for managing Docker containers**

![Docker Simple Panel](https://img.shields.io/badge/Docker-Simple_Panel-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Svelte](https://img.shields.io/badge/Svelte-5-FF3E00?style=for-the-badge&logo=svelte&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)

[![Quick Start](https://img.shields.io/badge/Quick_Start-Docker-2496ED?style=for-the-badge&logo=docker)](#quick-start-with-docker)
[![Documentation](https://img.shields.io/badge/Documentation-API-8A2BE2?style=for-the-badge)](#api-endpoints)

</div>

---

## 🌟 Overview

Docker Simple Panel is a sleek web application that provides a REST API and intuitive UI for managing Docker containers. Built with a Go backend and Svelte frontend, it offers authentication, real-time monitoring, and comprehensive container operations.

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🔐 **Authentication** | Single-user JWT-based authentication with configurable session timeout |
| 📦 **Container Management** | List, start, stop, restart containers with health status |
| 📊 **Real-time Logs** | WebSocket-powered log streaming with 30-minute history |
| 💾 **Volume Management** | View and manage Docker volumes with container associations |
| 📱 **PWA Support** | Install on iPhone, Android, and desktop with offline capabilities |
| ⚙️ **Configuration** | YAML-based config with web UI updates |
| 🐳 **Docker Integration** | Direct Docker daemon access via socket |

## 📋 Requirements

- **Go** 1.21+ (for building from source)
- **Docker** daemon running
- **Node.js** 18+ (for frontend development)

---

## 🚀 Quick Start with Docker

The easiest way to run Docker Simple Panel is using the pre-built Docker image:

```bash
# Pull the latest image
docker pull ghcr.io/dev-zapi/docker-simple-panel:latest

# Run the container
docker run -d \
  --name docker-simple-panel \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ./data:/app/data \
  ghcr.io/dev-zapi/docker-simple-panel:latest

# Verify it's running
curl http://localhost:8080/api/health
```

### 🔑 Default Credentials

| Username | Password |
|----------|----------|
| `admin` | `changeme` |

> ⚠️ **Important**: Change the default password immediately! Edit `/app/data/config.yaml` inside the container.

### 🔒 Security Best Practices

- 🛡️ Mount Docker socket only when necessary
- 🔑 Change default password on first use
- 🔐 Set a strong `JWT_SECRET` in production

---

## 💻 Development

### Dev Container (Recommended)

Use VS Code Dev Containers for a consistent development environment:

```bash
# Prerequisites: VS Code + Docker Desktop + Dev Containers extension
# Open project → Press F1 → "Dev Containers: Reopen in Container"
```

See [.devcontainer/README.md](.devcontainer/README.md) for details.

### Building from Source

```bash
# Clone repository
git clone https://github.com/dev-zapi/docker-simple-panel.git
cd docker-simple-panel

# Install dependencies
go mod download

# Build
go build -o docker-simple-panel .
```

### Building Docker Image Locally

```bash
docker build -t docker-simple-panel:local .

docker run -d \
  --name docker-simple-panel \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ./data:/app/data \
  docker-simple-panel:local
```

---

## ⚙️ Configuration

Configuration is managed via YAML file (default: `./config.yaml` or `CONFIG_PATH` env var).

### Example Configuration

```yaml
# Authentication
username: admin
password: changeme  # Auto-hashed with bcrypt

# Server settings
server:
  port: "8080"
  jwt_secret: "your-secret-key-change-in-production"
  session_max_timeout: 24  # hours

# Docker settings
docker:
  socket: "/var/run/docker.sock"
  volume_explorer_image: "ghcr.io/dev-zapi/docker-simple-panel:latest"

# Logging
logging:
  level: "info"  # error, warn, info, debug

# Static files (optional)
static_path: ""
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CONFIG_PATH` | `./config.yaml` | Path to YAML config file |
| `STATIC_PATH` | (empty) | Frontend static files directory |
| `SERVER_PORT` | `8080` | Server port |
| `JWT_SECRET` | (required in prod) | JWT signing secret |
| `DOCKER_SOCKET` | `/var/run/docker.sock` | Docker socket path |
| `DISABLE_REGISTRATION` | `false` | Disable user registration |
| `LOG_LEVEL` | `info` | Logging level |
| `VOLUME_EXPLORER_IMAGE` | (see above) | Volume explorer Docker image |

### 🔑 Password Management

- First startup: Creates default config with `admin`/`changeme`
- Password auto-hashed with bcrypt on save
- Change via web UI Settings page or edit config file
- **Critical**: Always change default password in production!

---

## ▶️ Running

```bash
# Default (uses ./config.yaml)
./docker-simple-panel

# Custom config
CONFIG_PATH=/path/to/config.yaml ./docker-simple-panel

# With static files
STATIC_PATH=/path/to/frontend/dist ./docker-simple-panel
```

---

## 📡 API Endpoints

Full OpenAPI v3 specification available at `GET /api/openapi.json`. Import into [Swagger UI](https://swagger.io/tools/swagger-ui/) for interactive documentation!

### Public Endpoints

#### Health Check
```http
GET /api/health
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "changeme"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "admin"
  }
}
```

### Protected Endpoints (Require JWT)

Include header: `Authorization: Bearer <token>`

#### List Containers
```http
GET /api/containers
```

Returns all containers with status and health info.

#### Container Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/containers/{id}` | Get container details |
| `POST` | `/api/containers/{id}/start` | Start container |
| `POST` | `/api/containers/{id}/stop` | Stop container |
| `POST` | `/api/containers/{id}/restart` | Restart container |
| `GET` | `/api/containers/{id}/logs/stream` | WebSocket log stream |

#### WebSocket Log Streaming

```http
GET /api/containers/{id}/logs/stream
```

Real-time log streaming with:
- ✅ JWT authentication (header or query param)
- ✅ 30-minute history
- ✅ Timestamps included
- ✅ stdout + stderr support

**Example:**
```javascript
const ws = new WebSocket(
  `ws://localhost:8080/api/containers/${containerId}/logs/stream`,
  [],
  { headers: { Authorization: `Bearer ${token}` } }
);

ws.onmessage = (event) => {
  console.log(event.data);
};
```

#### Docker Health
```http
GET /api/docker/health
```

#### Volume Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/volumes` | List volumes with containers |
| `DELETE` | `/api/volumes/{name}` | Delete volume |

#### Configuration

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/config` | Get current config |
| `PUT` | `/api/config` | Update config |

**Update Example:**
```http
PUT /api/config
Content-Type: application/json

{
  "docker_socket": "/custom/path/docker.sock",
  "log_level": "debug",
  "session_max_timeout": 48
}
```

---

## 🔒 Security

- 🔐 **Bcrypt** password hashing
- ⏰ JWT tokens expire (default: 24h)
- 🚫 Single-user authentication only
- ⚠️ Docker socket access grants full control
- 🔑 **Must change** default password & JWT secret in production

---

## 🧪 Development

### Run Tests
```bash
go test ./...
```

### Project Structure

```
.
├── main.go              # Entry point
├── config/              # YAML config management
├── database/            # SQLite operations
├── docker/              # Docker client wrapper
├── handlers/            # HTTP handlers
├── middleware/          # Auth, CORS, logging
├── models/              # Data models
├── webui/               # Svelte frontend
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   └── stores/
│   └── package.json
└── openapi.json         # API specification
```

---

## 🐳 Docker Image Building

### GitHub Actions (Automated)

1. Go to **Actions** tab
2. Select **"Build and Push Docker Image"**
3. Click **"Run workflow"**
4. Optionally specify custom tag
5. Builds multi-platform (amd64, arm64)

See [.github/workflows/README.md](.github/workflows/README.md) for details.

---

## 📄 License

MIT License - See LICENSE file for details.

---

<div align="center">

**Built with ❤️ using Go and Svelte**

[Report Issue](https://github.com/dev-zapi/docker-simple-panel/issues) • [Request Feature](https://github.com/dev-zapi/docker-simple-panel/issues) • [Discussions](https://github.com/dev-zapi/docker-simple-panel/discussions)

</div>
